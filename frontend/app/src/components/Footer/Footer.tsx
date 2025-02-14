import { View, Flex, Text } from "@adobe/react-spectrum";

export function Footer() {
  return (
    <View backgroundColor="gray-100" paddingX="size-100" paddingY="size-200">
      <Flex
        direction="row"
        alignItems="center"
        justifyContent="space-between"
        maxWidth={{ base: "100%", L: "1200px" }}
        marginX="auto"
      >
        <Text>© 2024 Neskuchka</Text>
      </Flex>
    </View>
  );
}
