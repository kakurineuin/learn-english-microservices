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
  FormErrorMessage,
  useDisclosure,
  VStack,
  FormControl,
  FormLabel,
  Input,
  Textarea,
  useToast,
  Switch,
} from '@chakra-ui/react';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as Yup from 'yup';
import { useAppDispatch, useAppSelector } from '../../store/hooks';
import { Exam } from '../../models/Exam';
import { updateExam } from '../../store/slices/examManagerSlice';

type FormData = {
  topic: string;
  description: string;
  isPublic: boolean;
};

type Props = {
  data: Exam;
};

function UpdateExamDialog({ data }: Props) {
  const user = useAppSelector((state) => state.session.user);
  const { isOpen, onOpen, onClose } = useDisclosure();

  const toast = useToast();
  const dispatch = useAppDispatch();

  const schema = Yup.object().shape({
    topic: Yup.string().trim().required().max(10),
    description: Yup.string().trim().required().max(100),
    isPublic: Yup.boolean().required(),
  });

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
    setFocus,
  } = useForm<FormData>({
    defaultValues: {
      topic: data.topic,
      description: data.description,
      isPublic: data.isPublic,
    },
    resolver: yupResolver(schema),
  });

  const onCloseHandler = () => {
    reset({
      topic: data.topic,
      description: data.description,
      isPublic: data.isPublic,
    });
    onClose();
  };

  const confirmOnClickHandler = handleSubmit(
    ({ topic, description, isPublic }) => {
      dispatch(
        updateExam({
          ...data,
          topic: topic.trim(),
          description: description.trim(),
          isPublic,
        }),
      )
        .unwrap()
        .then(() => {
          toast({
            title: '修改成功。',
            status: 'success',
            isClosable: true,
            position: 'top',
            variant: 'subtle',
          });

          data.topic = topic;
          data.description = description;
          data.isPublic = isPublic;
          onClose();
        })
        .catch((message) => {
          toast({
            title: '修改失敗。',
            description: message,
            status: 'error',
            isClosable: true,
            position: 'top',
            variant: 'subtle',
          });
        });
    },
  );

  const isPublicComponent =
    user!.role === 'admin' ? (
      <FormControl display="flex" alignItems="center">
        <FormLabel htmlFor="updateIsPublic" mb="0">
          Public
        </FormLabel>
        <Switch id="updateIsPublic" {...register('isPublic')} />
      </FormControl>
    ) : (
      <input
        id="isPublic"
        type="checkbox"
        style={{ display: 'none' }}
        {...register('isPublic')}
      />
    );

  useEffect(() => {
    if (isOpen) {
      setTimeout(() => {
        setFocus('topic');
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
          <ModalHeader>修改測驗</ModalHeader>
          <ModalCloseButton />

          <ModalBody>
            <form>
              <VStack spacing={3} align="start">
                <FormControl isInvalid={!!errors.topic}>
                  <FormLabel htmlFor="topic">Topic</FormLabel>
                  <Input
                    id="topic"
                    placeholder="Enter topic"
                    {...register('topic')}
                  />
                  <FormErrorMessage data-testid="topicError">
                    {errors.topic?.message}
                  </FormErrorMessage>
                </FormControl>

                <FormControl isInvalid={!!errors.description}>
                  <FormLabel htmlFor="description">Description</FormLabel>
                  <Textarea
                    id="description"
                    placeholder="Enter description"
                    rows={10}
                    {...register('description')}
                  />
                  <FormErrorMessage data-testid="descriptionError">
                    {errors.description?.message}
                  </FormErrorMessage>
                </FormControl>

                {isPublicComponent}
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

export default UpdateExamDialog;
