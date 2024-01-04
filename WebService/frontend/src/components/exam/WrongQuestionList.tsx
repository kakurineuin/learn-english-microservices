import { useMemo } from 'react';
import { Table, Thead, Tbody, Tr, Th, Td, Center, Box } from '@chakra-ui/react';
import { useTable, usePagination, Column, CellProps } from 'react-table';
import { v4 as uuidv4 } from 'uuid';
import { Question } from '../../models/Question';

type Props = {
  questions: Question[];
};

function WrongQuestionList({ questions }: Props) {
  const columns = useMemo(
    () =>
      [
        {
          Header: 'Ask',
          accessor: 'ask',
          Cell: ({ value }: CellProps<Question, string>) => <pre>{value}</pre>,
        },
        {
          Header: 'Answers',
          accessor: 'answers',
          Cell: ({ value }: CellProps<Question, string[]>) => {
            const answers = value;
            return answers.map((answer, i) => (
              <Box key={uuidv4()}>{`Answer${i + 1}: ${answer}`}</Box>
            ));
          },
        },
        {
          Header: 'Answer Wrong Times',
          accessor: 'answerWrongTimes',
          Cell: ({ value }: CellProps<Question, number>) => (
            <Center>{value}</Center>
          ),
        },
      ] as Column<Question>[],
    [],
  );
  const { getTableProps, getTableBodyProps, headerGroups, prepareRow, page } =
    useTable<Question>(
      {
        columns,
        data: questions,
        initialState: { pageIndex: 0 },
        manualPagination: true,
      },
      usePagination,
    );

  return (
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
  );
}

export default WrongQuestionList;
