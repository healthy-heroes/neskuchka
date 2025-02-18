import { Heading, View, Text, Avatar, Flex } from "@adobe/react-spectrum";

import { pageProps } from "../../constants";

export function TrackHeader() {
  return (
    <View
      paddingX="size-400"
      paddingY="size-200"
      backgroundColor="blue-200"
      colorVersion={6}
    >
      <View maxWidth={pageProps.maxWidth} marginX="auto">
        <Heading level={1}>Нескучный спорт</Heading>

        <Text>Тренируйтесь с нами — где бы вы ни находились!</Text>
        <br />
        <Text>
          Идеальная программа, чтобы поддерживать форму дома, без специального
          оборудования.
        </Text>

        <br />
        <br />
        <Flex alignItems="center" gap="size-40">
          <Avatar size="avatar-size-100" src="/avatar.jpg" alt="" />
          <Text>Илья Карягин</Text>
        </Flex>
      </View>
    </View>
  );
}
