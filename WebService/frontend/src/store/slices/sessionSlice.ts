import { PayloadAction, createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import axios, { AxiosError } from 'axios';
import { JwtPayload, jwtDecode } from 'jwt-decode';
import { loaderActions } from './loaderSlice';

type RequestParam = {
  username: string;
  password: string;
};

type User = {
  name: string;
  role: string;
  token: string;
};

type SessionState = {
  user: User | null;
};

type ExtendJwtPayload = {
  username: string;
  role: string;
};

type MyJwtPayload = ExtendJwtPayload & JwtPayload;

let user = null;
const token = localStorage.getItem('token');

if (token) {
  const { username: name, role } = jwtDecode<MyJwtPayload>(token);
  user = {
    name,
    role,
    token,
  };
}

const initialState: SessionState = {
  user,
};

export const signUp = createAsyncThunk(
  'signUp',
  async (requestParam: RequestParam, { dispatch, rejectWithValue }) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.post('/user', requestParam);
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

export const signIn = createAsyncThunk(
  'signIn',
  async (requestParam: RequestParam, { dispatch, rejectWithValue }) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.post('/login', requestParam);
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

export const sessionSlice = createSlice({
  name: 'session',
  initialState,
  reducers: {
    signOut: (state) => {
      state.user = null;
      localStorage.removeItem('token');
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(signUp.fulfilled, setUser)
      .addCase(signIn.fulfilled, setUser);
  },
});

function setUser(state: any, action: PayloadAction<User>) {
  const token = action.payload.token;
  localStorage.setItem('token', token);
  const decoded = jwtDecode<MyJwtPayload>(token);
  state.user = {
    name: decoded.username,
    role: decoded.role,
    token,
  };
}

export const sessionActions = sessionSlice.actions;
export const sessionReducer = sessionSlice.reducer;
