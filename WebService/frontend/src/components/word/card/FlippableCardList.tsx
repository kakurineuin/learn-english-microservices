import { useState } from 'react';
import {
  Flex,
  Tag,
  TagLabel,
  Text,
  Center,
  Box,
  Button,
  Spacer,
  Card,
  CardHeader,
  CardBody,
} from '@chakra-ui/react';
import { useHotkeys } from 'react-hotkeys-hook';
import { motion } from 'framer-motion';
import FlippableCard from './FlippableCard';
import { WordMeaning } from '../../../models/WordMeaning';
import AudioButton from '../../AudioButton';

const GAP = 100;
const CARD_WIDTH = 900;

type Props = {
  wordMeanings: WordMeaning[];
};

function FlippableCardList({ wordMeanings }: Props) {
  let initX: number;

  if (wordMeanings.length % 2 === 0) {
    initX =
      (CARD_WIDTH + GAP) / 2 +
      (wordMeanings.length / 2 - 1) * (CARD_WIDTH + GAP);
  } else {
    initX = Math.floor(wordMeanings.length / 2) * (CARD_WIDTH + GAP);
  }

  const [x, setX] = useState(initX);
  const [currentCardIndex, setCurrentCardIndex] = useState(0);
  const [isFlip, setIsFlip] = useState(false);
  const [isPlayUKAudio, setIsPlayUKAudio] = useState(false);
  const [isPlayUSAudio, setIsPlayUSAudio] = useState(false);
  const goToPreviousCard = () => {
    if (currentCardIndex === 0) {
      return;
    }

    setIsFlip(false);
    setX((prevX) => prevX + CARD_WIDTH + GAP);
    setCurrentCardIndex((prevCurrentCardIndex) => prevCurrentCardIndex - 1);
  };
  const goToNextCard = () => {
    if (currentCardIndex === wordMeanings.length - 1) {
      return;
    }

    setIsFlip(false);
    setX((prevX) => prevX - CARD_WIDTH - GAP);
    setCurrentCardIndex((prevCurrentCardIndex) => prevCurrentCardIndex + 1);
  };
  const flipCard = () => {
    // 若已經翻開就不處理
    if (isFlip) {
      return;
    }

    setIsFlip((prevIsFlip) => !prevIsFlip);
  };

  useHotkeys('a', goToPreviousCard, { scopes: ['card'] });
  useHotkeys('d', goToNextCard, { scopes: ['card'] });
  useHotkeys('s', flipCard, {
    scopes: ['card'],
  });

  useHotkeys('q', () => setIsPlayUKAudio(true), {
    scopes: ['card'],
  });

  useHotkeys('e', () => setIsPlayUSAudio(true), {
    scopes: ['card'],
  });

  return (
    <div
      style={{
        whiteSpace: 'nowrap',
        overflow: 'hidden',
      }}
    >
      <Center>
        <motion.div
          animate={{
            x,
          }}
          style={{
            marginTop: '0px',
            marginBottom: '50px',
          }}
        >
          {wordMeanings.map((wordMeaning, i) => {
            const {
              _id,
              partOfSpeech,
              gram,
              pronunciation: { text, ukAudioUrl, usAudioUrl },
              defGram,
              definition,
            } = wordMeaning;

            return (
              <div
                key={_id}
                style={{
                  display: 'inline-block',
                  marginLeft: i === 0 ? '0px' : `${GAP}px`,
                }}
              >
                <Card variant="outline" w="900px" mb="3">
                  <CardHeader>
                    <Flex alignItems="center">
                      {currentCardIndex === i && (
                        <Button
                          colorScheme="green"
                          variant="outline"
                          isDisabled={i === 0}
                          onClick={goToPreviousCard}
                        >
                          上一張
                        </Button>
                      )}

                      {partOfSpeech && (
                        <Tag
                          size="lg"
                          variant="outline"
                          colorScheme="blue"
                          ml="5"
                        >
                          <TagLabel>{partOfSpeech}</TagLabel>
                        </Tag>
                      )}

                      {gram && (
                        <Text color="blue.500" ml="2">
                          {gram}
                        </Text>
                      )}

                      <Flex alignItems="center" ml="10">
                        {text && <Text>{text}</Text>}

                        {ukAudioUrl &&
                          (i === currentCardIndex ? (
                            <AudioButton
                              colorScheme="orange"
                              variant="outline"
                              size="md"
                              ml="5"
                              audioUrl={ukAudioUrl}
                              isPlay={isPlayUKAudio}
                              onPlayFinished={() => setIsPlayUKAudio(false)}
                            >
                              UK
                            </AudioButton>
                          ) : (
                            <AudioButton
                              colorScheme="orange"
                              variant="outline"
                              size="md"
                              ml="5"
                              audioUrl={ukAudioUrl}
                            >
                              UK
                            </AudioButton>
                          ))}

                        {usAudioUrl &&
                          (i === currentCardIndex ? (
                            <AudioButton
                              colorScheme="orange"
                              variant="outline"
                              size="md"
                              ml="5"
                              audioUrl={usAudioUrl}
                              isPlay={isPlayUSAudio}
                              onPlayFinished={() => setIsPlayUSAudio(false)}
                            >
                              US
                            </AudioButton>
                          ) : (
                            <AudioButton
                              colorScheme="orange"
                              variant="outline"
                              size="md"
                              ml="5"
                              audioUrl={usAudioUrl}
                            >
                              US
                            </AudioButton>
                          ))}
                      </Flex>

                      <Spacer />

                      {currentCardIndex === i && (
                        <Button
                          colorScheme="green"
                          variant="outline"
                          isDisabled={i === wordMeanings.length - 1}
                          onClick={goToNextCard}
                        >
                          下一張
                        </Button>
                      )}
                    </Flex>
                  </CardHeader>

                  <CardBody overflowY="auto">
                    <Text color="blue.500">{defGram}</Text>
                    <Text fontSize="2xl" my="2" whiteSpace="pre-wrap">
                      {definition}
                    </Text>
                  </CardBody>
                </Card>

                <Box
                  onClick={() => {
                    // 若不是當前卡片，不用處理
                    if (currentCardIndex !== i) {
                      return;
                    }

                    flipCard();
                  }}
                >
                  <FlippableCard
                    isFlip={i === currentCardIndex ? isFlip : false}
                    wordMeaning={wordMeaning}
                  />
                </Box>
              </div>
            );
          })}
        </motion.div>
      </Center>
    </div>
  );
}

export default FlippableCardList;
