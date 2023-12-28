import { useEffect, useState } from 'react';
import {
  FormControl,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  Input,
  Kbd,
  Button,
  useToast,
  Badge,
  Tag,
  Box,
  TagLabel,
  Wrap,
  WrapItem,
} from '@chakra-ui/react';
import { Search2Icon } from '@chakra-ui/icons';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as Yup from 'yup';
import axios, { AxiosError } from 'axios';
import { loaderActions } from '../../store/slices/loaderSlice';
import { useAppDispatch } from '../../store/hooks';
import WordInfo from './WordInfo';
import { WordMeaning } from '../../models/WordMeaning';
import { Member } from '../../models/WordFamily';
import ShowText from '../ShowText';

type FormData = {
  word: string;
};

function WordForm() {
  const [wordFamilyMembers, setWordFamilyMembers] = useState<Member[]>([]);
  const [wordMeanings, setWordMeanings] = useState<WordMeaning[]>([]);
  const toast = useToast();
  const dispatch = useAppDispatch();

  const schema = Yup.object().shape({
    word: Yup.string().trim().required().max(50),
  });

  const {
    register,
    handleSubmit,
    formState: { errors },
    setFocus,
  } = useForm<FormData>({
    resolver: yupResolver(schema),
  });

  const searchWord = async (word: string) => {
    dispatch(loaderActions.toggleLoading());

    try {
      const response = await axios.get(`/restricted/word/${word.trim()}`);
      setWordFamilyMembers(
        response.data.wordFamily ? response.data.wordFamily.members : [],
      );
      setWordMeanings(response.data.wordMeanings);

      if (response.data.wordMeanings.length === 0) {
        toast({
          title: '查無資料。',
          description: '查不到此單字的解釋',
          status: 'info',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });
      } else {
        toast({
          title: '查詢成功。',
          description: '',
          status: 'success',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });
      }
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

  const submitHandler = handleSubmit(({ word }) => {
    searchWord(word);
  });

  let wordFamilyBox = null;

  if (wordFamilyMembers.length > 0) {
    wordFamilyBox = (
      <Box mt="5" borderWidth="1px" borderRadius="lg" p="3">
        <Badge colorScheme="purple">Word Family</Badge>
        <br />
        <Wrap spacing="5" align="center" mt="3">
          {wordFamilyMembers.map(({ partOfSpeech, word, _id }) => (
            <WrapItem alignItems="center" key={_id}>
              <Tag size="lg" variant="outline" colorScheme="blue" mr="2">
                <TagLabel>{partOfSpeech}</TagLabel>
              </Tag>
              <ShowText>{word}</ShowText>
            </WrapItem>
          ))}
        </Wrap>
      </Box>
    );
  }

  useEffect(() => {
    setFocus('word');
  }, [setFocus]);

  return (
    <>
      <form onSubmit={submitHandler}>
        <FormControl isInvalid={!!errors.word}>
          <FormLabel htmlFor="word">Word</FormLabel>
          <Input
            id="word"
            placeholder="Enter word"
            w="400px"
            {...register('word')}
          />
          <Button
            type="submit"
            leftIcon={<Search2Icon />}
            colorScheme="teal"
            variant="outline"
            ml="2"
          >
            Search
          </Button>
          <FormErrorMessage data-testid="wordError">
            {errors.word?.message}
          </FormErrorMessage>
          <FormHelperText>
            輸入英文單字後，點擊 [Search] 或按下<Kbd ml="1">enter</Kbd>
          </FormHelperText>
        </FormControl>
      </form>

      {wordFamilyBox}

      {wordMeanings.map(
        ({
          _id: wordMeaningId,
          word,
          partOfSpeech,
          gram,
          pronunciation,
          defGram,
          definition,
          examples,
          favoriteWordMeaningId,
        }) => (
          <WordInfo
            key={wordMeaningId!}
            wordMeaningId={wordMeaningId!}
            word={word}
            partOfSpeech={partOfSpeech}
            gram={gram}
            pronunciation={pronunciation}
            defGram={defGram}
            definition={definition}
            examples={examples}
            favoriteWordMeaningId={favoriteWordMeaningId!}
          />
        ),
      )}
    </>
  );
}

export default WordForm;
