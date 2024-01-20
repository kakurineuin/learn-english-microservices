import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import axios, { AxiosError } from 'axios';
import { loaderActions } from './loaderSlice';
import { ExamRecord } from '../../models/ExamRecord';

type ExamRecordManagerState = {
  total: number;
  pageCount: number;
  examRecords: ExamRecord[];
};

const initialState: ExamRecordManagerState = {
  total: 0,
  pageCount: 0,
  examRecords: [],
};

export const queryExamRecords = createAsyncThunk(
  'queryExamRecords',
  async (
    params: { pageIndex: number; pageSize: number; examId: string },
    { dispatch, rejectWithValue },
  ) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.get(
        `/restricted/exam/${params.examId}/record`,
        {
          params: {
            pageIndex: params.pageIndex,
            pageSize: params.pageSize,
          },
        },
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

type QueryExamRecords = {
  total: number;
  pageCount: number;
  examRecords: ExamRecord[];
};

export const examRecordManagerSlice = createSlice({
  name: 'examRecordManager',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(
      queryExamRecords.fulfilled,
      (state, action: PayloadAction<QueryExamRecords>) => {
        state.total = action.payload.total;
        state.pageCount = action.payload.pageCount;
        state.examRecords = action.payload.examRecords;
      },
    );
  },
});

export const examRecordManagerActions = examRecordManagerSlice.actions;
export const examRecordManagerReducer = examRecordManagerSlice.reducer;
