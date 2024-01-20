import { ReactNode } from 'react';
import { motion } from 'framer-motion';

type Props = {
  children: ReactNode;
  onAnimationComplete?: () => void;
};

function CardWithScale({ children, onAnimationComplete }: Props) {
  const cardVariants = {
    initial: {
      scale: 1,
    },
    show: {
      transition: {
        ease: 'easeOut',
        delay: 0.15,
        duration: 0.5,
      },
    },
    hover: {
      scale: 1.05,
    },
  };

  return (
    <motion.div
      variants={cardVariants}
      initial="initial"
      animate="show"
      whileHover="hover"
      // transition={{
      //   ease: 'easeOut',
      //   delay: 0.15,
      //   duration: 0.5,
      // }}
      className="w-96"
      onAnimationComplete={onAnimationComplete}
    >
      {children}
    </motion.div>
  );
}

export default CardWithScale;
