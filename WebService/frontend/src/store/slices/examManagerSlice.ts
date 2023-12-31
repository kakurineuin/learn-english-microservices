import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import axios, { AxiosError } from 'axios';
import { loaderActions } from './loaderSlice';
import { Exam } from '../../models/Exam';

type ExamManagerState = {
  saveExamTimes: number;
  total: number;
  pageCount: number;
  exams: Exam[];
};

const initialState: ExamManagerState = {
  saveExamTimes: 0,
  total: 0,
  pageCount: 0,
  exams: [],
};

export const queryExams = createAsyncThunk(
  'queryExams',
  async (
    params: { pageIndex: number; pageSize: number },
    { dispatch, rejectWithValue },
  ) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.get('/restricted/exam', { params });
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

export const createExam = createAsyncThunk(
  'createExam',
  async (exam: Exam, { dispatch, rejectWithValue }) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.post('/restricted/exam', exam);
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

export const updateExam = createAsyncThunk(
  'updateExam',
  async (exam: Exam, { dispatch, rejectWithValue }) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.patch('/restricted/exam', exam);
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

export const deleteExam = createAsyncThunk(
  'deleteExam',
  async (examId: string, { dispatch, rejectWithValue }) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.delete(`/restricted/exam/${examId}`);
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

type QueryExams = {
  total: number;
  pageCount: number;
  exams: Exam[];
};

export const examManagerSlice = createSlice({
  name: 'examManager',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(
        queryExams.fulfilled,
        (state, action: PayloadAction<QueryExams>) => {
          state.total = action.payload.total;
          state.pageCount = action.payload.pageCount;
          state.exams = action.payload.exams;
        },
      )
      .addCase(createExam.fulfilled, (state) => {
        state.saveExamTimes += 1;
      })
      .addCase(updateExam.fulfilled, (state) => {
        state.saveExamTimes += 1;
      })
      .addCase(deleteExam.fulfilled, (state) => {
        state.saveExamTimes += 1;
      });
  },
});

export const examManagerActions = examManagerSlice.actions;
export const examManagerReducer = examManagerSlice.reducer;
