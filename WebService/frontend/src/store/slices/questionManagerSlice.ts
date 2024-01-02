import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import axios, { AxiosError } from 'axios';
import { loaderActions } from './loaderSlice';
import { Question } from '../../models/Question';

type QuestionManagerState = {
  saveQuestionTimes: number;
  total: number;
  pageCount: number;
  questions: Question[];
};

const initialState: QuestionManagerState = {
  saveQuestionTimes: 0,
  total: 0,
  pageCount: 0,
  questions: [],
};

export const queryQuestions = createAsyncThunk(
  'queryQuestions',
  async (
    {
      examId,
      pageIndex,
      pageSize,
    }: { examId: string; pageIndex: number; pageSize: number },
    { dispatch, rejectWithValue },
  ) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.get(`/restricted/exam/${examId}/question`, {
        params: {
          pageIndex,
          pageSize,
        },
      });
      return response.data;
    } catch (err) {
      const errorMessage = axios.isAxiosError(err)
        ? (err as AxiosError<{ message: string }, any>).response!.data.message
        : '系統發生錯誤！';
      return rejectWithValue(errorMessage);
    } finally {
      dispatch(loaderActions.toggleLoading());
    }
  },
);

export const createQuestion = createAsyncThunk(
  'createQuestion',
  async (question: Question, { dispatch, rejectWithValue }) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.post(
        `/restricted/exam/${question.examId}/question`,
        question,
      );
      return response.data;
    } catch (err) {
      const errorMessage = axios.isAxiosError(err)
        ? (err as AxiosError<{ message: string }, any>).response!.data.message
        : '系統發生錯誤！';
      return rejectWithValue(errorMessage);
    } finally {
      dispatch(loaderActions.toggleLoading());
    }
  },
);

export const updateQuestion = createAsyncThunk(
  'updateQuestion',
  async (question: Question, { dispatch, rejectWithValue }) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.patch(
        `/restricted/exam/${question.examId}/question`,
        question,
      );
      return response.data;
    } catch (err) {
      const errorMessage = axios.isAxiosError(err)
        ? (err as AxiosError<{ message: string }, any>).response!.data.message
        : '系統發生錯誤！';
      return rejectWithValue(errorMessage);
    } finally {
      dispatch(loaderActions.toggleLoading());
    }
  },
);

export const deleteQuestion = createAsyncThunk(
  'deleteQuestion',
  async ({ examId, _id }: Question, { dispatch, rejectWithValue }) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.delete(
        `/restricted/exam/${examId}/question/${_id}`,
      );
      return response.data;
    } catch (err) {
      const errorMessage = axios.isAxiosError(err)
        ? (err as AxiosError<{ message: string }, any>).response!.data.message
        : '系統發生錯誤！';
      return rejectWithValue(errorMessage);
    } finally {
      dispatch(loaderActions.toggleLoading());
    }
  },
);

type QueryQuestions = {
  total: number;
  pageCount: number;
  questions: Question[];
};

export const questionManagerSlice = createSlice({
  name: 'questionManager',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(
        queryQuestions.fulfilled,
        (state, action: PayloadAction<QueryQuestions>) => {
          state.total = action.payload.total;
          state.pageCount = action.payload.pageCount;
          state.questions = action.payload.questions;
        },
      )
      .addCase(createQuestion.fulfilled, (state) => {
        state.saveQuestionTimes += 1;
      })
      .addCase(updateQuestion.fulfilled, (state) => {
        state.saveQuestionTimes += 1;
      })
      .addCase(deleteQuestion.fulfilled, (state) => {
        state.saveQuestionTimes += 1;
      });
  },
});

export const questionManagerActions = questionManagerSlice.actions;
export const questionManagerReducer = questionManagerSlice.reducer;
