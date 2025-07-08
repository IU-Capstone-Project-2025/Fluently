# backend/app.py

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import pandas as pd
from collections import Counter

# карта CEFR для внутренней логики (если понадобится)
CEFR_LEVELS = {'a1': 1, 'a2': 2, 'b1': 3, 'b2': 4, 'c1': 5, 'c2': 6}

class VocabularyRecommender:
    def __init__(self, dataset: pd.DataFrame):
        self.dataset = dataset
        self.dataset['cefr_num'] = self.dataset['CEFR_level'].str.lower().map(CEFR_LEVELS)
        # Normalize Total for frequency scoring (scale to [0, 1])
        self.dataset['normalized_total'] = (
            self.dataset['Total'] - self.dataset['Total'].min()
        ) / (self.dataset['Total'].max() - self.dataset['Total'].min())
        self.words_set = set(self.dataset['word'])

    def get_recommendations(self, known_words: list,
                            num_recommendations: int=10,
                            max_cefr_level: str='c2',
                            max_per_subtopic: int=2):
        max_cefr_num = CEFR_LEVELS[max_cefr_level.lower()]
        available_words = self.dataset[
            (~self.dataset['word'].isin(known_words)) &
            (self.dataset['cefr_num'] <= max_cefr_num)
        ].copy()

        if available_words.empty:
            return []

        known_words_df = self.dataset[self.dataset['word'].isin(known_words)]
        known_topics = set(known_words_df['topic'])
        known_subtopics = set(known_words_df['subtopic'])
        known_subsubtopics = set(known_words_df['subsubtopic'])

        # Compute hierarchical thematic scores
        def compute_theme_score(row):
            score = 0
            if row['topic'] in known_topics:
                score += 0.5
                if row['subtopic'] in known_subtopics:
                    score += 0.3
                    if row['subsubtopic'] in known_subsubtopics:
                        score += 0.2
            return score

        available_words['theme_score'] = available_words.apply(compute_theme_score, axis=1)


        available_words['final_score'] = (
            0.7 * available_words['theme_score'] +
            0.3 * available_words['normalized_total']
        )

        available_words = available_words.sort_values(
            by=['theme_score', 'normalized_total'],
            ascending=[False, False]
        )

        subtopic_counts = Counter()
        recommendations = []

        for _, row in available_words.iterrows():
            subtopic = row['subtopic']
            if subtopic_counts[subtopic] < max_per_subtopic:
                recommendations.append({
                    'word': row['word'],
                    'topic': row['topic'],
                    'subtopic': row['subtopic'],
                    'subsubtopic': row['subsubtopic'],
                    'CEFR_level': row['CEFR_level'],
                    'score': row['final_score']
                })
                subtopic_counts[subtopic] += 1
            if len(recommendations) >= num_recommendations:
                break

        return recommendations

# ——————————————————————————————————————————

# Инициализация
df = pd.read_csv('result.csv')
recommender = VocabularyRecommender(df)

app = FastAPI(title="Thesaurus API")

# — модели запросов
class RecommendRequest(BaseModel):
    words: list[str]

class HealthRequest(BaseModel):
    ping: str

# — health check
@app.post("/health")
async def health(req: HealthRequest):
    return {"status": "ok"}

# — endpoint рекомендаций
@app.post("/api/recommend")
async def recommend(req: RecommendRequest):
    # results = []
    # for w in req.words:
    #     recs = recommender.get_recommendations(w)
    #     # склейка всех рекомендаций в один список
    #     results.extend(recs)
    results = recommender.get_recommendations(req.words)
    return results

from fastapi.middleware.cors import CORSMiddleware

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],  # frontend origin
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)