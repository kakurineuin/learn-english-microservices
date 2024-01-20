export interface Question {
  _id?: string;
  examId: string;
  ask: string;
  answers: string[];
  userId?: string;
}
