import { QueryClient } from "@tanstack/react-query";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // данные считаются свежими 5 минут
      retry: 1, // одна повторная попытка при ошибке
    },
  },
});

export default queryClient;
