export type Member = {
  _id?: string;
  partOfSpeech: string;
  word: string;
};

export interface WordFamily {
  _id?: string;
  members: Member[];
}
