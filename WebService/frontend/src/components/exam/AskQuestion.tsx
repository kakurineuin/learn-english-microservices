import { useEffect, useState } from 'react';
import {
  Badge,
  Box,
  Flex,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Input,
} from '@chakra-ui/react';
import {
  useForm,
  useFieldArray,
  FieldErrorsImpl,
  DeepRequired,
} from 'react-hook-form';
import { motion, Variants } from 'framer-motion';
import { useAppDispatch } from '../../store/hooks';
import { Question } from '../../models/Question';
import { askFormActions } from '../../store/slices/askFormSlice';
import ShowText from '../ShowText';

type Props = {
  questionNumber: number;
  question: Question;
  isValidate: boolean;
  motionVariants: Variants;
};

type FormData = {
  answerObjects: {
    answer: string;
  }[];
};

function AskQuestion({
  questionNumber,
  question: { _id, ask, answers },
  isValidate,
  motionVariants,
}: Props) {
  const defaultValues: FormData = {
    answerObjects: [],
  };

  for (let i = 0; i < answers.length; i += 1) {
    defaultValues.answerObjects.push({
      answer: '',
    });
  }

  const {
    register,
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({
    defaultValues,
  });
  const { fields } = useFieldArray({
    name: 'answerObjects',
    control,
  });

  const dispatch = useAppDispatch();
  const [isShowAsk, setIsShowAsk] = useState(false);

  useEffect(() => {
    if (!isValidate) {
      return;
    }

    const submitHandler = handleSubmit(
      // eslint-disable-next-line
      (data: FormData) => {
        dispatch(
          askFormActions.addQuestionResult({
            questionId: _id!,
            isSuccess: true,
          }),
        );
      },

      // eslint-disable-next-line
      (formErrors: FieldErrorsImpl<DeepRequired<FormData>>) => {
        dispatch(
          askFormActions.addQuestionResult({
            questionId: _id!,
            isSuccess: false,
          }),
        );
      },
    );
    submitHandler();
  }, [_id, dispatch, handleSubmit, isValidate]);

  return (
    <motion.div
      variants={motionVariants}
      initial="hidden"
      animate="show"
      style={{ width: '100%' }}
      onAnimationComplete={() => setIsShowAsk(true)}
    >
      <Box borderWidth="1px" borderRadius="lg" w="100%">
        <form
          onSubmit={(e) => {
            e.preventDefault();
          }}
        >
          <Box m={3}>
            <Flex>
              <Badge my="2">Ask</Badge>
            </Flex>
            {isShowAsk && (
              <ShowText data-testid={`question${questionNumber}Ask`}>
                {ask}
              </ShowText>
            )}
          </Box>
          <Box m={3}>
            {fields.map((field, index) => {
              const label = `Answer${index + 1}`;
              const answer = answers[index];

              return (
                <FormControl
                  key={field.id}
                  isInvalid={!!errors?.answerObjects?.[index]?.answer}
                >
                  <FormLabel htmlFor={field.id}>
                    <Badge mt={2}>{label}</Badge>
                  </FormLabel>
                  <Input
                    id={field.id}
                    placeholder="Enter answer"
                    {...register(`answerObjects.${index}.answer` as const, {
                      validate: (value) =>
                        (value && value.trim() === answer) ||
                        `${label}: ${answer}`,
                    })}
                    isDisabled={isValidate}
                  />
                  <FormErrorMessage
                    data-testid={`question${questionNumber}Answer${
                      index + 1
                    }Error`}
                  >
                    {errors?.answerObjects?.[index]?.answer?.message}
                  </FormErrorMessage>
                </FormControl>
              );
            })}
          </Box>
        </form>
      </Box>
    </motion.div>
  );
}

export default AskQuestion;
