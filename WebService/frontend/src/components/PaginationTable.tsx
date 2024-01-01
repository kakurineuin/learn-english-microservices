import { useEffect, useState } from 'react';
import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Button,
  HStack,
  Text,
  Center,
} from '@chakra-ui/react';
import { useTable, usePagination, Column } from 'react-table';

type PaginationTableProps = {
  columns: Column<any>[];
  hiddenColumns?: string[];
  data: object[];
  fetchData: (pageIndex: number, pageSize: number) => void;
  pageCount: number;
  total: number;
  refreshId: number; // 當需要 refresh 時，就遞增此值
};

function PaginationTable({
  columns,
  hiddenColumns = [],
  data,
  fetchData,
  pageCount: controlledPageCount,
  total,
  refreshId,
}: PaginationTableProps) {
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    prepareRow,
    page,
    canPreviousPage,
    canNextPage,
    pageOptions,
    pageCount,
    gotoPage,
    nextPage,
    previousPage,
    state: { pageIndex, pageSize },
  } = useTable(
    {
      columns,
      data,
      initialState: { pageIndex: 0, hiddenColumns },
      manualPagination: true,
      pageCount: controlledPageCount,
    },
    usePagination,
  );

  const [prevRefreshId, setPrevRefreshId] = useState(refreshId);

  // TODO: 這是舊的寫法，測試 refreshId 功能沒問題後就刪除
  // // 重新查詢並顯示第1頁
  // useEffect(() => {
  //   if (pageIndex === 0) {
  //     fetchData(pageIndex, pageSize);
  //   } else {
  //     gotoPage(0);
  //   }
  //
  //   // eslint-disable-next-line
  // }, [refreshId]);

  useEffect(() => {
    // 重新查詢並顯示第1頁
    if (prevRefreshId !== refreshId) {
      if (pageIndex === 0) {
        fetchData(pageIndex, pageSize);
      } else {
        gotoPage(0);
      }

      setPrevRefreshId(refreshId);
    } else {
      fetchData(pageIndex, pageSize);
    }
  }, [fetchData, pageIndex, pageSize, refreshId, gotoPage, prevRefreshId]);

  return (
    <>
      <Table {...getTableProps()}>
        <Thead>
          {headerGroups.map((headerGroup) => (
            <Tr {...headerGroup.getHeaderGroupProps()}>
              {headerGroup.headers.map((column) => (
                <Th {...column.getHeaderProps()}>
                  <Center>{column.render('Header')}</Center>
                </Th>
              ))}
            </Tr>
          ))}
        </Thead>
        <Tbody {...getTableBodyProps()}>
          {page.map((row) => {
            prepareRow(row);
            return (
              <Tr {...row.getRowProps()}>
                {row.cells.map((cell) => (
                  <Td {...cell.getCellProps()}>{cell.render('Cell')}</Td>
                ))}
              </Tr>
            );
          })}
        </Tbody>
      </Table>

      {/* pagination */}
      <HStack mt="2" justify="space-between">
        <HStack spacing={2}>
          <Button
            colorScheme="teal"
            variant="outline"
            size="xs"
            onClick={() => gotoPage(0)}
            isDisabled={!canPreviousPage}
          >
            {'<<'}
          </Button>
          <Button
            colorScheme="teal"
            variant="outline"
            size="xs"
            onClick={() => previousPage()}
            isDisabled={!canPreviousPage}
          >
            {'<'}
          </Button>
          <Button
            colorScheme="teal"
            variant="outline"
            size="xs"
            onClick={() => nextPage()}
            isDisabled={!canNextPage}
          >
            {'>'}
          </Button>
          <Button
            colorScheme="teal"
            variant="outline"
            size="xs"
            onClick={() => gotoPage(pageCount - 1)}
            isDisabled={!canNextPage}
          >
            {'>>'}
          </Button>
          <span>
            Page{' '}
            <strong>
              {pageIndex + 1} of {pageOptions.length}
            </strong>
          </span>
        </HStack>
        <Text fontSize="md">
          每頁最多顯示 {pageSize} 筆，總共 {total} 筆
        </Text>
      </HStack>
    </>
  );
}

export default PaginationTable;
