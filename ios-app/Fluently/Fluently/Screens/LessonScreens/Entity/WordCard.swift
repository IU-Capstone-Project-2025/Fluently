//
//  WordCardExs.swift
//  Fluently
//
//  Created by Савва Пономарев on 25.06.2025.
//

import Foundation


final class WordCard: Exercise {
    var wordId: UUID
    var word: String
    var translation: String
    var transcription: String
    var cefrLevel: String
    var isNew: Bool
    var topic: String
    var subtopic: String
    var sentences: [Sentence]
    var exercise: Exercise

    init(
        exerciseId: UUID,
        wordId: UUID,
        word: String,
        translation: String,
        transcription: String,
        cefrLevel: String,
        isNew:Bool,
        topic: String,
        subtopic: String,
        sentences: [Sentence],
        exercise: Exercise
    ) {
        self.wordId = wordId
        self.word = word
        self.translation = translation
        self.transcription = transcription
        self.cefrLevel = cefrLevel
        self.isNew = isNew
        self.topic = topic
        self.subtopic = subtopic
        self.sentences = sentences
        self.exercise = exercise

        super.init(
            exerciseId: exerciseId,
            exerciseType: "wordCard",
            correctAnswer: word
        )
    }
}

// MARK: - Generator for mock lesson

struct WordCardGenerator {
    static func generateCards(count: Int = 5) -> [WordCard] {
        let sampleWords = [
            ("apple", "яблоко", "[ˈæpəl]", "A1", "Food", "Fruits",
             "I ate an apple for breakfast.", "Я съел яблоко на завтрак.",
             ["banana", "orange", "pear"]),

            ("house", "дом", "[haʊs]", "A1", "Home", "Buildings",
             "I built a house last year.", "Я построил дом в прошлом году.",
             ["building", "garage", "shed"]),

            ("car", "машина", "[kɑːr]", "A1", "Transport", "Vehicles",
             "He drives a car to work.", "Он ездит на машине на работу.",
             ["bus", "bike", "truck"]),

            ("book", "книга", "[bʊk]", "A1", "Education", "Objects",
             "She read a book yesterday.", "Она прочитала книгу вчера.",
             ["magazine", "newspaper", "notebook"]),

            ("dog", "собака", "[dɒɡ]", "A1", "Animals", "Pets",
             "We have a dog at home.", "У нас есть собака дома.",
             ["cat", "rabbit", "hamster"])
        ]

        return (0..<min(count, sampleWords.count)).map { index in
            let word = sampleWords[index]
            return createCard(
                word: word.0,
                translation: word.1,
                transcription: word.2,
                cefrLevel: word.3,
                topic: word.4,
                subtopic: word.5,
                sentence: word.6,
                sentenceTranslation: word.7,
                wrongOptions: word.8
            )
        }
    }

    private static func createCard(
        word: String,
        translation: String,
        transcription: String,
        cefrLevel: String,
        topic: String,
        subtopic: String,
        sentence: String,
        sentenceTranslation: String,
        wrongOptions: [String]
    ) -> WordCard {
        let wordId = UUID()
        let sentenceId = UUID()
        let exerciseId = UUID()

        // Create sentence with blank for the word
        let template = sentence.replacingOccurrences(of: word, with: "____")

        // Combine wrong options with correct answer and shuffle
        var allOptions = wrongOptions + [word]
        allOptions.shuffle()

        let type = [
            ExerciseType.pickOptions,
            ExerciseType.chooseTranslationEngRuss,
            ExerciseType.typeTranslationRussEng
        ].randomElement()

        return WordCard(
            exerciseId: UUID(),
            wordId: wordId,
            word: word,
            translation: translation,
            transcription: transcription,
            cefrLevel: cefrLevel,
            isNew: Bool.random(),
            topic: topic,
            subtopic: subtopic,
            sentences: [
                Sentence(
                    sentenceId: sentenceId,
                    sentece: sentence,
                    translation: sentenceTranslation
                )
            ],

            exercise: {
                switch type {
                    case .chooseTranslationEngRuss:
                        return ChooseTranslationExs (
                            exerciseId: UUID(),
                            wordId: wordId,
                            word: translation,
                            options: allOptions,
                            correctAnswer: word
                        )
                    case .typeTranslationRussEng:
                        return TypeTranslationExs(
                            exerciseId: exerciseId,
                            wordId: wordId,
                            word: word,
                            correctAnswer: translation
                        )
                    case .pickOptions:
                        return PickOptionsExs(
                            exerciseId: exerciseId,
                            sentenceId: sentenceId,
                            sentence: template,
                            options: allOptions,
                            correctAnswer: word
                        )
                    case .recordPronounce:
                        return ChooseTranslationExs (
                            exerciseId: exerciseId,
                            wordId: wordId,
                            word: translation,
                            options: allOptions,
                            correctAnswer: word
                        )
                    case .wordCard:
                        return PickOptionsExs(
                            exerciseId: exerciseId,
                            sentenceId: sentenceId,
                            sentence: template,
                            options: allOptions,
                            correctAnswer: word
                        )
                    case .numberOfWords:
                        return TypeTranslationExs(
                            exerciseId: exerciseId,
                            wordId: wordId,
                            word: word,
                            correctAnswer: translation
                        )
                    case nil:
                        return PickOptionsExs(
                            exerciseId: exerciseId,
                            sentenceId: sentenceId,
                            sentence: template,
                            options: allOptions,
                            correctAnswer: word
                        )
                }
            }()
        )
    }
}
