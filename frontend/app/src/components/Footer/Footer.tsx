import { View, Flex, Text } from "@adobe/react-spectrum";

import { pageProps } from "../../pages/constants";

export function Footer() {
  return (
    <View backgroundColor="gray-100" paddingX="size-100" paddingY="size-200">
      <Flex
        direction="row"
        alignItems="center"
        justifyContent="space-between"
        maxWidth={pageProps.maxWidth}
        marginX="auto"
      >
        <Text>© 2024 Neskuchka</Text>
      </Flex>
    </View>
  );
}
