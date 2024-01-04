export interface AnswerWrong {
  _id?: string;
  examId: string;
  questionId: string;
  times: number;
  userId: string;
}
