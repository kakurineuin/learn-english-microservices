import { useEffect } from 'react';
import {
  Container,
  FormControl,
  FormLabel,
  FormErrorMessage,
  Input,
  Button,
  VStack,
  useToast,
} from '@chakra-ui/react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as Yup from 'yup';
import PageHeading from '../components/PageHeading';
import { useAppDispatch } from '../store/hooks';
import { signIn } from '../store/slices/sessionSlice';

type FormData = {
  username: string;
  password: string;
};

export default function SignIn() {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const toast = useToast();

  const schema = Yup.object().shape({
    username: Yup.string().trim().required().max(20),
    password: Yup.string().trim().required().min(8).max(20),
  });

  const {
    register,
    handleSubmit,
    formState: { errors },
    setFocus,
  } = useForm<FormData>({
    resolver: yupResolver(schema),
  });

  const submitHandler = handleSubmit(({ username, password }) => {
    dispatch(
      signIn({
        username,
        password,
      }),
    )
      .unwrap()
      .then(() => {
        toast({
          title: '登入成功',
          status: 'success',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });

        // 轉向首頁
        navigate('/');
      })
      .catch((message) => {
        toast({
          title: '登入失敗',
          description: message,
          status: 'error',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });
      });
  });

  useEffect(() => {
    setFocus('username');
  }, [setFocus]);

  return (
    <Container>
      <PageHeading title="登入" />
      <form onSubmit={submitHandler}>
        <VStack spacing={3}>
          <FormControl isInvalid={!!errors.username}>
            <FormLabel htmlFor="username">Username</FormLabel>
            <Input
              id="username"
              placeholder="Enter username"
              {...register('username')}
            />
            <FormErrorMessage>{errors.username?.message}</FormErrorMessage>
          </FormControl>

          <FormControl isInvalid={!!errors.password}>
            <FormLabel htmlFor="password">Password</FormLabel>
            <Input
              id="password"
              type="password"
              placeholder="Enter password"
              {...register('password')}
            />
            <FormErrorMessage>{errors.password?.message}</FormErrorMessage>
          </FormControl>

          <Button colorScheme="blue" variant="outline" type="submit">
            確定
          </Button>
        </VStack>
      </form>
    </Container>
  );
}
