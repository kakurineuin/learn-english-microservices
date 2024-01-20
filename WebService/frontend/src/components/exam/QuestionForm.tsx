import { useEffect } from 'react';
import {
  VStack,
  FormControl,
  FormLabel,
  FormErrorMessage,
  Textarea,
  Input,
  Button,
  useToast,
} from '@chakra-ui/react';
import { useForm, useFieldArray } from 'react-hook-form';
import { AddIcon, DeleteIcon } from '@chakra-ui/icons';
import { useAppDispatch } from '../../store/hooks';
import { createQuestion } from '../../store/slices/questionManagerSlice';

type FormData = {
  ask: string;
  answers: { answer: string }[];
};

type Props = {
  examId: string;
};

const askMaxLength = 500;
const answerMaxLength = 30;

function QuestionForm({ examId }: Props) {
  const toast = useToast();
  const dispatch = useAppDispatch();
  const defaultValues: FormData = {
    ask: '',
    answers: [{ answer: '' }],
  };
  const {
    register,
    control,
    handleSubmit,
    reset,
    formState: { errors },
    setFocus,
  } = useForm<FormData>({
    defaultValues,
  });
  const { fields, append, remove } = useFieldArray({
    name: 'answers',
    control,
  });
  const submitHandler = handleSubmit(({ ask, answers }) => {
    dispatch(
      createQuestion({
        examId,
        ask: ask.trim(),
        answers: answers.map((item) => item.answer.trim()),
      }),
    )
      .unwrap()
      .then(() => {
        toast({
          title: '新增成功',
          status: 'success',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });

        // 清除表單
        reset();
      })
      .catch((message) => {
        toast({
          title: '新增失敗',
          description: message,
          status: 'error',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });
      });
  });

  useEffect(() => {
    setFocus('ask');
  }, [setFocus]);

  return (
    <form onSubmit={submitHandler}>
      <VStack spacing={3} align="start">
        <FormControl isInvalid={!!errors.ask}>
          <FormLabel htmlFor="ask">Ask</FormLabel>
          <Textarea
            id="ask"
            placeholder="Enter ask"
            {...register('ask', {
              required: 'ask is a required field',
              maxLength: {
                value: askMaxLength,
                message: `ask must be at most ${askMaxLength} characters`,
              },
            })}
          />
          <FormErrorMessage data-testid="askError">
            {errors.ask?.message}
          </FormErrorMessage>
        </FormControl>
        {fields.map((field, index) => {
          const inputId = `answer${index + 1}`;
          const label = `Answer${index + 1}`;
          return (
            <FormControl
              isInvalid={!!errors?.answers?.[index]?.answer}
              key={field.id}
            >
              <FormLabel htmlFor={inputId}>{label}</FormLabel>
              <span>
                <Input
                  id={inputId}
                  placeholder="Enter answer"
                  {...register(`answers.${index}.answer` as const, {
                    required: `${label} is a required field`,
                    maxLength: {
                      value: answerMaxLength,
                      message: `${label} must be at most ${answerMaxLength} characters`,
                    },
                  })}
                  w="200px"
                />
                {index === 0 ? (
                  <Button
                    colorScheme="teal"
                    variant="outline"
                    ml="5"
                    onClick={() => append({ answer: '' })}
                    data-testid="addAnswer"
                  >
                    <AddIcon />
                  </Button>
                ) : (
                  <Button
                    colorScheme="red"
                    variant="outline"
                    ml="5"
                    onClick={() => remove(index)}
                    data-testid={`deleteAnswer${index + 1}`}
                  >
                    <DeleteIcon />
                  </Button>
                )}
              </span>
              <FormErrorMessage data-testid={`answer${index + 1}Error`}>
                {errors?.answers?.[index]?.answer?.message}
              </FormErrorMessage>
            </FormControl>
          );
        })}
        <Button colorScheme="blue" variant="outline" type="submit">
          新增
        </Button>
      </VStack>
    </form>
  );
}

export default QuestionForm;
