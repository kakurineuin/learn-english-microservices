import {
  Center,
  Container,
  ListItem,
  Text,
  UnorderedList,
  useToast,
} from '@chakra-ui/react';
import axios, { AxiosError } from 'axios';
import { useEffect, useState } from 'react';
import PageHeading from '../components/PageHeading';
import FlippableCardList from '../components/word/card/FlippableCardList';
import { WordMeaning } from '../models/WordMeaning';
import ShowText from '../components/ShowText';
import { useAppDispatch } from '../store/hooks';
import { loaderActions } from '../store/slices/loaderSlice';

function WordCard() {
  const [favoriteWordMeanings, setFavoriteWordMeanings] = useState<
    WordMeaning[]
  >([]);
  const dispatch = useAppDispatch();
  const toast = useToast();

  useEffect(() => {
    const queryRandomFavoriteWordMeanings = async () => {
      dispatch(loaderActions.toggleLoading());

      try {
        const response = await axios.get('/restricted/word/card');
        setFavoriteWordMeanings(response.data.favoriteWordMeanings);
      } catch (err) {
        const errorMessage = axios.isAxiosError(err)
          ? (err as AxiosError<{ message: string }, any>).response!.data.message
          : '系統發生錯誤！';
        toast({
          title: '查詢失敗。',
          description: errorMessage,
          status: 'error',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });
      } finally {
        dispatch(loaderActions.toggleLoading());
      }
    };

    queryRandomFavoriteWordMeanings();
  }, [dispatch, toast]);

  return (
    <>
      <Container maxW="container.xl" mt="3">
        <PageHeading title="單字卡">
          <ShowText fontSize="lg">
            此頁會從最愛的單字解釋中隨機產生卡片用來複習
          </ShowText>
          <br />
          <ShowText as="b" fontSize="lg">
            滑鼠操作
          </ShowText>
          <UnorderedList>
            <ListItem>
              <ShowText fontSize="lg">點擊卡片即可翻面</ShowText>
            </ListItem>
            <ListItem>
              <ShowText fontSize="lg">點擊 [下一張] 移到下一張卡片</ShowText>
            </ListItem>
            <ListItem>
              <ShowText fontSize="lg">點擊 [上一張] 移到上一張卡片</ShowText>
            </ListItem>
          </UnorderedList>
          <br />
          <ShowText as="b" fontSize="lg">
            鍵盤操作
          </ShowText>
          <UnorderedList>
            <ListItem>
              <ShowText fontSize="lg">按 S 鍵即可翻面</ShowText>
            </ListItem>
            <ListItem>
              <ShowText fontSize="lg">按 D 鍵移到下一張卡片</ShowText>
            </ListItem>
            <ListItem>
              <ShowText fontSize="lg">按 A 鍵移到上一張卡片</ShowText>
            </ListItem>
          </UnorderedList>
        </PageHeading>
      </Container>
      {favoriteWordMeanings && favoriteWordMeanings.length > 0 ? (
        <FlippableCardList wordMeanings={favoriteWordMeanings} />
      ) : (
        <Center>
          <Text fontSize="xl">
            請先使用[查詢單字]功能查詢單字，將喜歡的單字解釋加入最愛，再回到這裡就會有單字卡了。
          </Text>
        </Center>
      )}
    </>
  );
}

export default WordCard;
