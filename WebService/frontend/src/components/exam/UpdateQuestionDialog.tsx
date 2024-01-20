import { useEffect } from 'react';
import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Button,
  useDisclosure,
  VStack,
  FormControl,
  FormLabel,
  FormErrorMessage,
  Input,
  Textarea,
  useToast,
} from '@chakra-ui/react';
import { AddIcon, DeleteIcon } from '@chakra-ui/icons';
import { useForm, useFieldArray } from 'react-hook-form';
import { useAppDispatch } from '../../store/hooks';
import { Question } from '../../models/Question';
import { updateQuestion } from '../../store/slices/questionManagerSlice';

type FormData = {
  ask: string;
  answers: { answer: string }[];
};

type Props = {
  data: Question;
};

const askMaxLength = 500;
const answerMaxLength = 30;

function UpdateQuestionDialog({ data }: Props) {
  const { isOpen, onOpen, onClose } = useDisclosure();

  const toast = useToast();
  const dispatch = useAppDispatch();

  const {
    register,
    control,
    handleSubmit,
    reset,
    formState: { errors },
    setFocus,
  } = useForm<FormData>({
    defaultValues: {
      ask: data.ask,
      answers: data.answers.map((answer) => ({ answer })),
    },
  });
  const { fields, append, remove } = useFieldArray({
    name: 'answers',
    control,
  });

  const onCloseHandler = () => {
    reset({
      ask: data.ask,
      answers: data.answers.map((answer) => ({ answer })),
    });
    onClose();
  };

  const confirmOnClickHandler = handleSubmit(({ ask, answers }) => {
    dispatch(
      updateQuestion({
        _id: data._id,
        examId: data.examId,
        ask: ask.trim(),
        answers: answers.map((item) => item.answer.trim()),
      })
    )
      .unwrap()
      .then(() => {
        toast({
          title: '修改成功',
          status: 'success',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });

        data.ask = ask;
        data.answers = answers.map((item) => item.answer.trim());
        onClose();
      })
      .catch((message) => {
        toast({
          title: '修改失敗',
          description: message,
          status: 'error',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });
      });
  });

  useEffect(() => {
    if (isOpen) {
      setTimeout(() => {
        setFocus('ask');
      }, 100);
    }
  }, [isOpen, setFocus]);

  return (
    <>
      <Button colorScheme="teal" variant="outline" size="sm" onClick={onOpen}>
        修改
      </Button>

      <Modal size="5xl" isOpen={isOpen} onClose={onCloseHandler}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>修改題目</ModalHeader>
          <ModalCloseButton />

          <ModalBody>
            <form>
              <VStack spacing={3} align="start">
                <FormControl isInvalid={!!errors.ask}>
                  <FormLabel htmlFor="ask">Ask</FormLabel>
                  <Textarea
                    id="ask"
                    placeholder="Enter ask"
                    rows={10}
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
              </VStack>
            </form>
          </ModalBody>

          <ModalFooter>
            <Button
              colorScheme="blue"
              variant="outline"
              mr={3}
              onClick={confirmOnClickHandler}
            >
              Confirm
            </Button>
            <Button variant="outline" onClick={onCloseHandler}>
              Cancel
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
}

export default UpdateQuestionDialog;
