import { useEffect } from 'react';
import {
  VStack,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Textarea,
  Input,
  Button,
  useToast,
  Switch,
} from '@chakra-ui/react';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as Yup from 'yup';
import { useAppDispatch, useAppSelector } from '../../store/hooks';
import { createExam } from '../../store/slices/examManagerSlice';

type FormData = {
  topic: string;
  description: string;
  isPublic: boolean;
};

function ExamForm() {
  const user = useAppSelector((state) => state.session.user);
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
    resolver: yupResolver(schema),
  });

  const submitHandler = handleSubmit(({ topic, description, isPublic }) => {
    dispatch(
      createExam({
        topic: topic.trim(),
        description: description.trim(),
        tags: [],
        isPublic,
      }),
    )
      .unwrap()
      .then(() => {
        toast({
          title: '新增成功。',
          status: 'success',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });

        // 清除表單
        reset({ topic: '', description: '' });
      })
      .catch((message) => {
        toast({
          title: '新增失敗。',
          description: message,
          status: 'error',
          isClosable: true,
          position: 'top',
          variant: 'subtle',
        });
      });
  });

  const isPublicComponent =
    user!.role === 'admin' ? (
      <FormControl display="flex" alignItems="center">
        <FormLabel htmlFor="isPublic" mb="0">
          Public
        </FormLabel>
        <Switch id="isPublic" {...register('isPublic')} />
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
    setFocus('topic');
  }, [setFocus]);

  return (
    <form onSubmit={submitHandler}>
      <VStack spacing={3} align="start">
        <FormControl isInvalid={!!errors.topic}>
          <FormLabel htmlFor="topic">Topic</FormLabel>
          <Input id="topic" placeholder="Enter topic" {...register('topic')} />
          <FormErrorMessage data-testid="topicError">
            {errors.topic?.message}
          </FormErrorMessage>
        </FormControl>

        <FormControl isInvalid={!!errors.description}>
          <FormLabel htmlFor="description">Description</FormLabel>
          <Textarea
            id="description"
            placeholder="Enter description"
            {...register('description')}
          />
          <FormErrorMessage data-testid="descriptionError">
            {errors.description?.message}
          </FormErrorMessage>
        </FormControl>

        {isPublicComponent}

        <Button colorScheme="blue" variant="outline" type="submit">
          新增
        </Button>
      </VStack>
    </form>
  );
}

export default ExamForm;
