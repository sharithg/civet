import { MutationFunction, useMutation } from "@tanstack/react-query";
import _ from "lodash";
import { useCallback } from "react";

export function useDebouncedMutation<TData, TVariables>(input: {
  mutationFn: MutationFunction<TData, TVariables>;
}) {
  const mutation = useMutation({ mutationFn: input.mutationFn });

  const debouncedMutate = useCallback(
    _.debounce((data) => {
      mutation.mutate(data);
    }, 500),
    []
  );

  return { ...mutation, debouncedMutate };
}
