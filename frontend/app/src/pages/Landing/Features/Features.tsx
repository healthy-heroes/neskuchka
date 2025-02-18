import { motion } from "framer-motion";
import {
  Content,
  Grid,
  View,
  Heading,
  Text,
  Flex,
} from "@adobe/react-spectrum";

import { pageProps } from "../../constants";

import styles from "./Features.module.css";

const features = [
  {
    title: "Для тренеров",
    description:
      "Создавайте программы тренировок и следите за прогрессом учеников",
    icon: "🎯",
  },
  {
    title: "Для учеников",
    description:
      "Получайте персонализированные тренировки и отмечайте результаты",
    icon: "💪",
  },
  {
    title: "Прогресс",
    description: "Отслеживайте достижения и анализируйте свои результаты",
    icon: "📈",
  },
];

export function Features() {
  const container = {
    hidden: { opacity: 0 },
    show: {
      opacity: 1,
      transition: { staggerChildren: 0.1 },
    },
  };

  const item = {
    hidden: { opacity: 0, y: 20 },
    show: { opacity: 1, y: 0 },
  };

  return (
    <View paddingX={pageProps.paddingX} paddingY="size-1000">
      <Content>
        <Flex
          direction="column"
          gap="size-400"
          maxWidth={pageProps.maxWidth}
          marginX="auto"
        >
          <motion.div
            variants={container}
            initial="hidden"
            whileInView="show"
            viewport={{ once: true }}
          >
            <Grid columns={{ base: "1fr", M: "repeat(3, 1fr)" }} gap="size-400">
              {features.map((feature) => (
                <motion.div key={feature.title} variants={item}>
                  <View
                    borderRadius="medium"
                    borderWidth="thin"
                    borderColor="gray-200"
                    padding="size-400"
                    UNSAFE_className={styles.featureCard}
                  >
                    <Flex direction="column">
                      <Text UNSAFE_className={styles.icon}>{feature.icon}</Text>
                      <Heading
                        level={2}
                        marginBottom="size-100"
                        UNSAFE_className={styles.featureTitle}
                      >
                        {feature.title}
                      </Heading>
                      <Text UNSAFE_className={styles.featureDescription}>
                        {feature.description}
                      </Text>
                    </Flex>
                  </View>
                </motion.div>
              ))}
            </Grid>
          </motion.div>
        </Flex>
      </Content>
    </View>
  );
}
