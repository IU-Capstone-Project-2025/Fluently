# Data

This is one the only places the `ui` layer can communicate with.
This README file describes conceptually each of the repositories.

## `LessonRepository`
Currently, the server supports only single operation: generating the lesson for the user. Therefore, it seems reasonable to use the following model:

On a user device, there is a single memory cell for storing the ongoing lesson and updating its progress. What can we do with this memory cell?
- Fetch and replace its content with the response from the server
- Update the current lesson component.
- Observe the current `LessonComponent` as `Flow`.
- Clear the cell.

That's it. Currently, `LessonRepository` and its methods follow precicely this model.



