export type Pronunciation = {
  text: string;
  ukAudioUrl: string;
  usAudioUrl: string;
};

export type Sentence = {
  audioUrl: string;
  text: string;
};

export type Example = {
  pattern: string;
  examples: Sentence[];
};

export interface WordMeaning {
  _id?: string;
  word: string;
  partOfSpeech: string;
  gram: string;
  pronunciation: Pronunciation;
  defGram: string;
  definition: string;
  examples: Example[];
  orderByNo: number;
  queryByWords: string[];
  favoriteWordMeaningId?: string;
}
