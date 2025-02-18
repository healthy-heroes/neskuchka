import { Heading, Text, View } from "@adobe/react-spectrum";

import { pageProps } from "../../constants";

export function Workouts() {
  return (
    <View
      paddingX="size-400"
      paddingY="size-400"
      maxWidth={pageProps.maxWidth}
      marginX="auto"
    >
      <View
        borderWidth="thin"
        borderColor="gray-200"
        borderRadius="medium"
        paddingX="size-200"
        paddingBottom="size-200"
      >
        <Heading level={2}>Тренировка от 17 февраля</Heading>

        <Heading level={3}>Разминка</Heading>
        <Text>
          3 раунда:
          <br />
          - 10 подвижность 90-90
          <br />
          - 5+5 скалолаз+уход в дельфина на одной
          <br />- 20 мельница в наклоне
        </Text>

        <Heading level={3}>Комплекс</Heading>
        <Text>
          5 ранудов не на время:
          <br />
          - 10+10 болгарские выпады
          <br />
          - 10+10 боковая планка на колене
          <br />- 20 пресс на прямых руки над головой
        </Text>
      </View>

      <View
        marginTop="size-400"
        borderWidth="thin"
        borderColor="gray-200"
        borderRadius="medium"
        paddingX="size-200"
        paddingBottom="size-200"
      >
        <Heading level={2}>Тренировка от 24 февраля</Heading>

        <Heading level={3}>Разминка</Heading>
        <Text>
          3 раунда:
          <br />
          - 5+5+5+5 вращение бедра стоя
          <br />
          - 10 пресс на прямых
          <br />- 10 раскрытий в седе
        </Text>

        <Heading level={3}>Комплекс</Heading>
        <Text>
          5 ранудов не на время:
          <br />
          - 50 пингвиньих прыжков
          <br />
          - 10 обратных берпи
          <br />- 10+10 выпрыгиваний в выпаде без смены
        </Text>
      </View>
    </View>
  );
}
