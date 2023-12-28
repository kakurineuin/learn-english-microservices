import { useCallback, useEffect, useState } from 'react';
import { Button, useToast } from '@chakra-ui/react';
import { AiOutlineSound } from 'react-icons/ai';

type ColorScheme =
  | 'whiteAlpha'
  | 'blackAlpha'
  | 'gray'
  | 'red'
  | 'orange'
  | 'yellow'
  | 'green'
  | 'teal'
  | 'blue'
  | 'cyan'
  | 'purple'
  | 'pink'
  | 'linkedin'
  | 'facebook'
  | 'messenger'
  | 'whatsapp'
  | 'twitter'
  | 'telegram';
type Variant = 'ghost' | 'outline' | 'solid' | 'link' | 'unstyled';
type Size = 'lg' | 'md' | 'sm' | 'xs';

type Props = {
  colorScheme: ColorScheme;
  variant: Variant;
  size: Size;
  audioUrl: string;
  isPlay?: boolean;
  onPlayFinished?: () => void;
  children?: React.ReactNode;
  mt?: string | number;
  mb?: string | number;
  ml?: string | number;
  mr?: string | number;
};

function AudioButton({
  colorScheme,
  variant,
  size,
  audioUrl,
  isPlay,
  onPlayFinished,
  children,
  mt,
  mb,
  ml,
  mr,
}: Props) {
  const [isAudioLoading, setIsAudioLoading] = useState(false);
  const toast = useToast();

  const playAudio = useCallback(async () => {
    setIsAudioLoading(true);

    try {
      const audio = new Audio(audioUrl);
      await audio.play();
    } catch (err) {
      console.error('Failed to play audio, error', err);
      toast({
        title: '抱歉，此音檔不可使用。',
        description: '',
        status: 'error',
        isClosable: true,
        position: 'top',
        variant: 'subtle',
      });
    } finally {
      setIsAudioLoading(false);

      if (onPlayFinished) {
        onPlayFinished();
      }
    }
  }, [audioUrl, toast, onPlayFinished]);

  useEffect(() => {
    if (!isPlay) {
      return;
    }

    playAudio();
  }, [playAudio, isPlay]);

  return (
    <Button
      isLoading={isAudioLoading}
      leftIcon={<AiOutlineSound />}
      colorScheme={colorScheme}
      variant={variant}
      size={size}
      onClick={() => playAudio()}
      mt={mt}
      mb={mb}
      ml={ml}
      mr={mr}
    >
      {children}
    </Button>
  );
}

export default AudioButton;
