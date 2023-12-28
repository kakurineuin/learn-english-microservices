import {
  Flex,
  Box,
  Heading,
  Button,
  Spacer,
  Text,
  Center,
  useToast,
} from '@chakra-ui/react';
import { useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { sessionActions } from '../store/slices/sessionSlice';

function Menus() {
  const dispatch = useAppDispatch();
  const user = useAppSelector((state) => state.session.user);
  const navigate = useNavigate();
  const toast = useToast();

  const homeClickHandler = () => {
    navigate('/');
  };
  // const examClickHandler = () => {
  //   navigate('/restricted/exam');
  // };
  const wordClickHandler = () => {
    navigate('/restricted/word');
  };
  const favoriteWordMeaningClickHandler = () => {
    navigate('/restricted/word/favorite');
  };
  const WordCardClickHandler = () => {
    navigate('/restricted/word/card');
  };
  const signUpHandler = () => {
    navigate('/signup');
  };
  const signInHandler = () => {
    navigate('/signin');
  };
  const signOutHandler = () => {
    dispatch(sessionActions.signOut());
    toast({
      title: '登出成功',
      status: 'success',
      isClosable: true,
      position: 'top',
      variant: 'subtle',
    });
  };

  return (
    <Flex p="3" className="w-full backdrop-blur-md bg-gray-950/60">
      <Box p="2">
        <Heading size="md">Learn English</Heading>
      </Box>
      <Spacer />
      {user && (
        <Center mr={8}>
          <Text fontSize="xl">{user.name}</Text>
        </Center>
      )}
      <Box>
        {!user && (
          <Button
            colorScheme="red"
            variant="outline"
            mr="2"
            onClick={signUpHandler}
          >
            註冊
          </Button>
        )}

        {!user && (
          <Button
            colorScheme="teal"
            variant="outline"
            mr="2"
            onClick={signInHandler}
          >
            登入
          </Button>
        )}

        {user && (
          <Button
            colorScheme="teal"
            variant="outline"
            mr="2"
            onClick={signOutHandler}
          >
            登出
          </Button>
        )}

        <Button
          colorScheme="teal"
          variant="outline"
          mr="2"
          onClick={homeClickHandler}
        >
          首頁
        </Button>

        {
          // TODO: 補上
          // <Button
          //   colorScheme="teal"
          //   variant="outline"
          //   mr="2"
          //   isDisabled={!username}
          //   onClick={examClickHandler}
          // >
          //   測驗管理 {!username && '(請先登入)'}
          // </Button>
        }

        <Button
          colorScheme="teal"
          variant="outline"
          mr="2"
          isDisabled={!user}
          onClick={wordClickHandler}
        >
          查詢單字 {!user && '(請先登入)'}
        </Button>

        <Button
          colorScheme="teal"
          variant="outline"
          mr="2"
          isDisabled={!user}
          onClick={favoriteWordMeaningClickHandler}
        >
          最愛的單字解釋 {!user && '(請先登入)'}
        </Button>

        <Button
          colorScheme="teal"
          variant="outline"
          mr="2"
          isDisabled={!user}
          onClick={WordCardClickHandler}
        >
          單字卡 {!user && '(請先登入)'}
        </Button>
      </Box>
    </Flex>
  );
}

export default Menus;
