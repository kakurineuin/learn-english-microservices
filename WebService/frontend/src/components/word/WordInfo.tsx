import { useState } from 'react';
import {
  Button,
  Box,
  Flex,
  Heading,
  Tag,
  TagLabel,
  UnorderedList,
  ListItem,
  Divider,
  Spacer,
  useToast,
} from '@chakra-ui/react';
import { MdOutlineFavoriteBorder, MdOutlineFavorite } from 'react-icons/md';
import axios, { AxiosError } from 'axios';
import { loaderActions } from '../../store/slices/loaderSlice';
import { useAppDispatch } from '../../store/hooks';
import { Pronunciation, Example } from '../../models/WordMeaning';
import AudioButton from '../AudioButton';
import ShowText from '../ShowText';

type Props = {
  wordMeaningId: string;
  word: string;
  partOfSpeech: string;
  gram: string;
  pronunciation: Pronunciation;
  defGram: string;
  definition: string;
  examples: Example[];
  favoriteWordMeaningId: string;
};

function WordInfo({
  wordMeaningId,
  word,
  partOfSpeech,
  gram,
  pronunciation: { text: pronText, ukAudioUrl, usAudioUrl },
  defGram,
  definition,
  examples,
  favoriteWordMeaningId,
}: Props) {
  const toast = useToast();
  const dispatch = useAppDispatch();
  const [favoriteId, setFavoriteId] = useState(favoriteWordMeaningId);

  const exampleComponents = examples.map(
    ({ pattern, examples: exampleArray }) => {
      const components = exampleArray.map(({ audioUrl, text }) => (
        <ListItem mt="2" key={text}>
          <Flex alignItems="center">
            <AudioButton
              colorScheme="teal"
              variant="outline"
              size="sm"
              mr="2"
              audioUrl={audioUrl}
            />
            <ShowText>{text}</ShowText>
          </Flex>
        </ListItem>
      ));

      if (pattern) {
        return (
          <Box my="7" key={pattern}>
            <ShowText as="b" fontSize="xl">
              {pattern}
            </ShowText>
            <UnorderedList>{components}</UnorderedList>
          </Box>
        );
      }

      return components[0];
    },
  );

  const createFavoriteWordMeaning = async () => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.post('/word/favorite', {
        wordMeaningId,
      });
      setFavoriteId(response.data.favoriteWordMeaningId);
      toast({
        title: '新增最愛的單字解釋成功。',
        status: 'success',
        isClosable: true,
        position: 'top',
        variant: 'subtle',
      });
    } catch (err) {
      const errorMessage = axios.isAxiosError(err)
        ? (err as AxiosError<{ message: string }, any>).response!.data.message
        : '系統發生錯誤！';
      toast({
        title: '新增最愛的單字解釋失敗。',
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

  const deleteFavoriteWordMeaning = async () => {
    dispatch(loaderActions.toggleLoading());

    try {
      await axios.delete(`/word/favorite/${favoriteId}`);
      setFavoriteId('');
      toast({
        title: '刪除最愛的單字解釋成功。',
        status: 'success',
        isClosable: true,
        position: 'top',
        variant: 'subtle',
      });
    } catch (err) {
      const errorMessage = axios.isAxiosError(err)
        ? (err as AxiosError<{ message: string }, any>).response!.data.message
        : '系統發生錯誤！';
      toast({
        title: '刪除最愛的單字解釋失敗。',
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

  return (
    <Box mt="5" borderWidth="1px" borderRadius="lg" p="3">
      <Box>
        <Flex alignItems="center">
          <Heading as="h2" size="xl">
            {word}
          </Heading>

          {partOfSpeech && (
            <Tag size="lg" variant="outline" colorScheme="blue" ml="5">
              <TagLabel>{partOfSpeech}</TagLabel>
            </Tag>
          )}

          {gram && (
            <ShowText color="blue.500" ml="2">
              {gram}
            </ShowText>
          )}

          <Flex alignItems="center" ml="10">
            {pronText && <ShowText>{pronText}</ShowText>}

            {ukAudioUrl && (
              <AudioButton
                colorScheme="orange"
                variant="outline"
                size="md"
                ml="5"
                audioUrl={ukAudioUrl}
              >
                UK
              </AudioButton>
            )}

            {usAudioUrl && (
              <AudioButton
                colorScheme="orange"
                variant="outline"
                size="md"
                ml="5"
                audioUrl={usAudioUrl}
              >
                US
              </AudioButton>
            )}
          </Flex>

          <Spacer />

          <Button
            leftIcon={
              favoriteId ? <MdOutlineFavorite /> : <MdOutlineFavoriteBorder />
            }
            colorScheme={favoriteId ? 'red' : 'gray.500'}
            variant="outline"
            onClick={() =>
              favoriteId
                ? deleteFavoriteWordMeaning()
                : createFavoriteWordMeaning()
            }
          >
            Favorite
          </Button>
        </Flex>
      </Box>

      <Box>
        <ShowText color="blue.500" mb="2">
          {defGram}
        </ShowText>
        <ShowText fontSize="2xl" key={definition}>
          {definition}
        </ShowText>
        <Divider mt="2" />
        <UnorderedList mt="5">{exampleComponents}</UnorderedList>
      </Box>
    </Box>
  );
}

export default WordInfo;
