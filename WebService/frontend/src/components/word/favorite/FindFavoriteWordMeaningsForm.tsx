import { useEffect } from 'react';
import {
  FormControl,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  Input,
  Kbd,
  Button,
} from '@chakra-ui/react';
import { Search2Icon } from '@chakra-ui/icons';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as Yup from 'yup';

type Props = {
  onSubmit: (word: string) => void;
  onReset: () => void;
};

type FormData = {
  word: string;
};

function FindFavoriteWordMeaningsForm({ onSubmit, onReset }: Props) {
  const schema = Yup.object().shape({
    word: Yup.string().trim().required().max(50),
  });

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
    setFocus,
  } = useForm<FormData>({
    resolver: yupResolver(schema),
  });

  const submitHandler = handleSubmit(({ word }) => {
    onSubmit(word);
  });
  const resetHandler = () => {
    reset();
    onReset();
  };

  useEffect(() => {
    setFocus('word');
  }, [setFocus]);

  return (
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

        <Button variant="outline" ml="2" onClick={resetHandler}>
          Reset
        </Button>

        <FormErrorMessage data-testid="wordError">
          {errors.word?.message}
        </FormErrorMessage>
        <FormHelperText>
          輸入英文單字後，點擊 [Search] 或按下<Kbd ml="1">enter</Kbd>
        </FormHelperText>
      </FormControl>
    </form>
  );
}

export default FindFavoriteWordMeaningsForm;
