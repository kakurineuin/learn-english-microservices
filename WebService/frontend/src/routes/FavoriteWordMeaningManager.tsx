import { useMemo, useCallback, useState } from 'react';
import {
  Box,
  Container,
  UnorderedList,
  Heading,
  ListItem,
  VStack,
  useToast,
  Tag,
  TagLabel,
  Divider,
  Flex,
} from '@chakra-ui/react';
import { CellProps, Column } from 'react-table';
import axios, { AxiosError } from 'axios';
import { loaderActions } from '../store/slices/loaderSlice';
import { useAppDispatch } from '../store/hooks';
import { Pronunciation, WordMeaning } from '../models/WordMeaning';
import PaginationTable from '../components/PaginationTable';
import ConfirmDialog from '../components/ConfirmDialog';
import AudioButton from '../components/AudioButton';
import FindFavoriteWordMeaningsForm from '../components/word/favorite/FindFavoriteWordMeaningsForm';
import PageHeading from '../components/PageHeading';
import ShowText from '../components/ShowText';

function FavoriteWordMeaningManager() {
  const dispatch = useAppDispatch();
  const [refreshId, setRefreshId] = useState(0);
  const toast = useToast();
  const [word, setWord] = useState('');

  const columns = useMemo(
    () =>
      [
        {
          Header: 'Word',
          accessor: 'word',
          Cell: ({ value }: CellProps<WordMeaning, string>) => (
            <Heading as="h4" size="md">
              {value}
            </Heading>
          ),
        },
        {
          Header: 'Gram',
          accessor: 'gram',
        },
        {
          Header: 'Part Of Speech',
          accessor: 'partOfSpeech',
          Cell: ({ value, row }: CellProps<WordMeaning, string>) => (
            <div>
              {
                // 朗文字典有些單字解釋是沒有加上詞性的
                value && (
                  <Tag size="lg" variant="outline" colorScheme="blue" ml="5">
                    <TagLabel>{value}</TagLabel>
                  </Tag>
                )
              }
              {row.values.gram && (
                <ShowText key={row.values.gram} color="blue.500" mt="2">
                  {row.values.gram}
                </ShowText>
              )}
            </div>
          ),
        },
        {
          Header: 'Pronunciation',
          accessor: 'pronunciation',
          Cell: ({
            value: { text, ukAudioUrl, usAudioUrl },
          }: CellProps<WordMeaning, Pronunciation>) => (
            <VStack>
              <ShowText key={text}>{text}</ShowText>

              {ukAudioUrl && (
                <AudioButton
                  colorScheme="orange"
                  variant="outline"
                  size="sm"
                  audioUrl={ukAudioUrl}
                >
                  UK
                </AudioButton>
              )}

              {usAudioUrl && (
                <AudioButton
                  colorScheme="orange"
                  variant="outline"
                  size="sm"
                  audioUrl={usAudioUrl}
                >
                  US
                </AudioButton>
              )}
            </VStack>
          ),
        },
        {
          Header: 'defGram',
          accessor: 'defGram',
        },
        {
          Header: 'examples',
          accessor: 'examples',
        },
        {
          Header: 'Definition',
          accessor: 'definition',
          Cell: ({ value, row }: CellProps<WordMeaning, string>) => (
            <div>
              {row.values.defGram && (
                <ShowText key={row.values.defGram} color="blue.500">
                  {row.values.defGram}
                </ShowText>
              )}

              <ShowText key={value} fontSize="2xl" my="3">
                {value}
              </ShowText>

              <Divider />

              <UnorderedList mt="5">
                {(row.values as WordMeaning).examples.map(
                  ({ pattern, examples: exampleArray }) => {
                    const components = exampleArray.map(
                      ({ audioUrl, text }) => (
                        <ListItem mt="2" key={text}>
                          <Flex alignItems="center">
                            <AudioButton
                              colorScheme="teal"
                              variant="outline"
                              size="sm"
                              audioUrl={audioUrl}
                            />
                            <ShowText key={text} ml="2">
                              {text}
                            </ShowText>
                          </Flex>
                        </ListItem>
                      ),
                    );

                    if (pattern) {
                      return (
                        <Box my="7" key={pattern}>
                          <ShowText key={pattern} as="b" fontSize="xl">
                            {pattern}
                          </ShowText>
                          <UnorderedList>{components}</UnorderedList>
                        </Box>
                      );
                    }

                    return components[0];
                  },
                )}
              </UnorderedList>
            </div>
          ),
        },
        {
          Header: 'Edit',
          accessor: '_id',
          width: 50,
          Cell: ({ row }: CellProps<WordMeaning, string>) => (
            <VStack spacing={2}>
              <ConfirmDialog
                button={{
                  colorScheme: 'red',
                  variant: 'outline',
                  size: 'sm',
                  text: '刪除',
                }}
                header="刪除確認"
                yesCallback={() => {
                  const deleteFavoriteWordMeaning = async (
                    favoriteId: string,
                  ) => {
                    dispatch(loaderActions.toggleLoading());

                    try {
                      await axios.delete(
                        `/restricted/word/favorite/${favoriteId}`,
                      );

                      // Refresh table
                      setRefreshId((prevRefreshId) => prevRefreshId + 1);

                      toast({
                        title: '刪除最愛的單字解釋成功。',
                        status: 'success',
                        isClosable: true,
                        position: 'top',
                        variant: 'subtle',
                      });
                    } catch (err) {
                      const errorMessage = axios.isAxiosError(err)
                        ? (err as AxiosError<{ message: string }, any>)
                            .response!.data.message
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

                  deleteFavoriteWordMeaning(
                    row.original.favoriteWordMeaningId!,
                  );
                }}
              >
                確定要刪除此單字解釋？
              </ConfirmDialog>
            </VStack>
          ),
        },
      ] as Column<WordMeaning>[],
    [dispatch, toast],
  );

  const hiddenColumns = useMemo(() => ['gram', 'defGram', 'examples'], []);
  const [total, setTotal] = useState(0);
  const [pageCount, setPageCount] = useState(0);
  const [data, setData] = useState<WordMeaning[]>([]);
  const fetchData = useCallback(
    async (pageIndex: number, pageSize: number) => {
      dispatch(loaderActions.toggleLoading());

      try {
        const response = await axios.get('/restricted/word/favorite', {
          params: {
            word,
            pageIndex,
            pageSize,
          },
        });
        setData(response.data.favoriteWordMeanings);
        setTotal(response.data.total);
        setPageCount(response.data.pageCount);
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
    },
    [dispatch, toast, word],
  );

  const submitHandler = (queryWord: string) => {
    setWord(queryWord);

    // 重設為第1頁
    setRefreshId((prevRefreshId) => prevRefreshId + 1);
  };
  const resetHandler = () => {
    submitHandler('');
  };

  return (
    <Container maxW="container.xl" mt="3">
      <PageHeading title="最愛的單字解釋">
        <ShowText fontSize="lg">
          此頁列出最愛的單字解釋，資料來自朗文線上字典(https://www.ldoceonline.com/)
        </ShowText>
        <UnorderedList>
          <ListItem>
            <ShowText fontSize="lg">
              查詢最愛的單字解釋 - 查詢出已被加入最愛的單字解釋
            </ShowText>
          </ListItem>
          <ListItem>
            <ShowText fontSize="lg">
              刪除最愛的單字解釋 - 點擊[刪除]會顯示確認對話框，點擊[Yes]即可刪除
            </ShowText>
          </ListItem>
          <ListItem>
            <ShowText fontSize="lg">
              若要新增最愛的單字解釋，請使用選單上的[查詢單字]功能
            </ShowText>
          </ListItem>
        </UnorderedList>
      </PageHeading>
      <FindFavoriteWordMeaningsForm
        onSubmit={submitHandler}
        onReset={resetHandler}
      />
      <Divider my="10" />
      <PaginationTable
        columns={columns}
        hiddenColumns={hiddenColumns}
        data={data}
        fetchData={fetchData}
        pageCount={pageCount}
        total={total}
        refreshId={refreshId}
      />
    </Container>
  );
}

export default FavoriteWordMeaningManager;
