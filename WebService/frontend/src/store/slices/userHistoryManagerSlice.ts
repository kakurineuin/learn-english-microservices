import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import axios, { AxiosError } from 'axios';
import { loaderActions } from './loaderSlice';
import { UserHistory } from '../../models/UserHistory';

type UserHistoryManagerState = {
  total: number;
  pageCount: number;
  userHistories: UserHistory[];
};

const initialState: UserHistoryManagerState = {
  total: 0,
  pageCount: 0,
  userHistories: [],
};

export const queryUserHistories = createAsyncThunk(
  'queryUserHistories',
  async (
    params: { pageIndex: number; pageSize: number },
    { dispatch, rejectWithValue },
  ) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.get('/restricted/user/history', { params });
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

type QueryUserHistories = {
  total: number;
  pageCount: number;
  userHistories: UserHistory[];
};

export const userHistoryManagerSlice = createSlice({
  name: 'userHistoryManager',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(
      queryUserHistories.fulfilled,
      (state, action: PayloadAction<QueryUserHistories>) => {
        state.total = action.payload.total;
        state.pageCount = action.payload.pageCount;
        state.userHistories = action.payload.userHistories;
      },
    );
  },
});

export const userHistoryManagerActions = userHistoryManagerSlice.actions;
export const userHistoryManagerReducer = userHistoryManagerSlice.reducer;
