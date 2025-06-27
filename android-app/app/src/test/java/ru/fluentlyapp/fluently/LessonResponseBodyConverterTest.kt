package ru.fluentlyapp.fluently

import org.junit.Test
import ru.fluentlyapp.fluently.network.convertToLesson
import ru.fluentlyapp.fluently.testing.mockLessonResponse
import org.junit.Assert.*
import org.junit.Before

class LessonResponseBodyConverterTest {
    val convertedLesson = mockLessonResponse.convertToLesson()


    @Before
    fun setup() {
        println(convertedLesson)
    }

    @Test
    fun converterProducesCorrectNumberOfComponents() {
        assertEquals(5, convertedLesson.components.size)
    }
}
