// import { useState, useEffect } from 'react';
import {
  Container,
  // SimpleGrid,
  // useToast,
  UnorderedList,
  ListItem,
} from '@chakra-ui/react';
// import axios, { AxiosError } from 'axios';
// import { motion } from 'framer-motion';
// import { useAppDispatch } from '../store/hooks';
// import ExamInfoCard from '../components/exam/ExamInfoCard';
// import { loaderActions } from '../store/slices/loaderSlice';
import ShowText from '../components/ShowText';
import PageHeading from '../components/PageHeading';

function Home() {
  // const [examInfos, setExamInfos] = useState([]);
  // const dispatch = useAppDispatch();
  // const toast = useToast();

  // TODO: 之後恢復
  // useEffect(() => {
  //   const queryExamInfos = async () => {
  //     dispatch(loaderActions.toggleLoading());
  //
  //     try {
  //       const response = await axios.get('/exam/info');
  //       setExamInfos(response.data.examInfos);
  //     } catch (err) {
  //       const errorMessage = axios.isAxiosError(err)
  //         ? (err as AxiosError<{ message: string }, any>).response!.data.message
  //         : '系統發生錯誤！';
  //       toast({
  //         title: '查詢失敗。',
  //         description: errorMessage,
  //         status: 'error',
  //         isClosable: true,
  //         position: 'top',
  //         variant: 'subtle',
  //       });
  //     } finally {
  //       dispatch(loaderActions.toggleLoading());
  //     }
  //   };
  //
  //   queryExamInfos();
  // }, [dispatch, toast]);

  return (
    <Container maxW="container.xl">
      <PageHeading title="測驗集">
        <ShowText fontSize="lg">此頁會列出可以使用的測驗</ShowText>
        <UnorderedList>
          <ListItem>
            <ShowText fontSize="lg">點擊[開始測驗]進行測驗</ShowText>
          </ListItem>
          <ListItem>
            <ShowText fontSize="lg">點擊[成績紀錄]查看紀錄</ShowText>
          </ListItem>
        </UnorderedList>
        <br />
        <ShowText fontSize="lg">選單說明</ShowText>
        <UnorderedList>
          <ListItem>
            <ShowText fontSize="lg">
              [測驗管理] - 可以新增、修改、刪除自己的測驗
            </ShowText>
          </ListItem>
          <ListItem>
            <ShowText fontSize="lg">
              [查詢單字] - 可以查詢英文單字的解釋
            </ShowText>
          </ListItem>
          <ListItem>
            <ShowText fontSize="lg">
              [最愛的單字解釋] - 列出從[查詢單字]加入的最愛的單字解釋
            </ShowText>
          </ListItem>
          <ListItem>
            <ShowText fontSize="lg">
              [單字卡] - 從最愛的單字解釋隨機產生單字卡片以供複習
            </ShowText>
          </ListItem>
        </UnorderedList>
      </PageHeading>

      <UnorderedList>
        <ListItem>
          <ShowText fontSize="lg">
            請先註冊或登入，才能使用全部功能，註冊只需要帳號和密碼即可，不需要
            email
          </ShowText>
        </ListItem>
        <ListItem>
          <ShowText fontSize="lg">點擊標題旁邊的驚嘆號會顯示功能說明</ShowText>
        </ListItem>
      </UnorderedList>

      {
        // TODO: 之後恢復
        // <SimpleGrid
        //   mt="5"
        //   minChildWidth="400px"
        //   spacing="20px"
        //   data-testid="grid"
        // >
        //   {examInfos.map(
        //     (
        //       {
        //         examId,
        //         topic,
        //         description,
        //         isPublic,
        //         questionCount,
        //         recordCount,
        //       },
        //       index,
        //     ) => (
        //       <motion.div
        //         initial={{ opacity: 0, scale: 0.5 }}
        //         animate={{ opacity: 1, scale: 1 }}
        //         transition={{ duration: 0.5, delay: 0.2 * index }}
        //         key={examId}
        //       >
        //         <ExamInfoCard
        //           examId={examId}
        //           topic={topic}
        //           description={description}
        //           isPublic={isPublic}
        //           questionCount={questionCount}
        //           recordCount={recordCount}
        //         />
        //       </motion.div>
        //     ),
        //   )}
        // </SimpleGrid>
      }
    </Container>
  );
}

export default Home;
