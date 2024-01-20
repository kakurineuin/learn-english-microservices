export interface Exam {
  _id?: string;
  topic: string;
  description: string;
  tags: string[];
  isPublic: boolean;
  userId?: string;
}
