import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

type QuestionResult = {
  questionId: string;
  isSuccess: boolean;
};

type AskFormState = {
  questionResults: QuestionResult[];
};

const initialState: AskFormState = {
  questionResults: [],
};

export const askFormSlice = createSlice({
  name: 'askForm',
  initialState,
  reducers: {
    addQuestionResult: (state, action: PayloadAction<QuestionResult>) => {
      state.questionResults.push(action.payload);
    },
  },
});

export const askFormActions = askFormSlice.actions;
export const askFormReducer = askFormSlice.reducer;
