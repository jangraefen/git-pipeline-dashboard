import { Box, Skeleton, VStack } from "@hope-ui/solid"
import { createReaction, createResource, createSignal, onCleanup, Show, For } from "solid-js"
import Pipeline from "./Pipeline"

const pipelineFetcher = async ({ repositoryID, repositorySource }) =>
  (await fetch(`/repositories/${repositorySource}/${repositoryID}`)).json()

export default (props) => {
  const [pipelines, { refetch: refetchPipelines }] = createResource(() => ({ ...props }), pipelineFetcher)
  const interval = setInterval(refetchPipelines, 1 * 60 * 1000)
  onCleanup(() => clearInterval(interval))

  const [initiallyLoaded, setInitiallyLoaded] = createSignal(false)
  const loadTracker = createReaction(() => setInitiallyLoaded(true))
  loadTracker(pipelines)

  return (
    <Box>
      <Skeleton loaded={initiallyLoaded}>
        <Show when={pipelines.error == null && pipelines() != null && pipelines().length > 0 && !pipelines().status}>
          <VStack spacing="$2">
            <For each={pipelines()}>{(pipeline) => <Pipeline pipeline={pipeline} />}</For>
          </VStack>
        </Show>
        <Show when={pipelines.error == null && pipelines() != null && !!pipelines().status}>
          <Box bg="$neutral7" borderRadius="$md" p="$3" fontStyle="italic" fontSize="$xs">
            <Box mt="$2">
              <Box fontSize="$xs">Error while fetching pipelines:</Box>
              <Box fontSize="$xs">
                {pipelines().state} - {pipelines().title}
              </Box>
              <Box fontSize="$xs">{pipelines().detail}</Box>
            </Box>
          </Box>
        </Show>
        <Show when={pipelines.error == null && pipelines() != null && pipelines().length == 0}>
          <Box bg="$neutral7" borderRadius="$md" p="$3" fontStyle="italic" fontSize="$xs">
            No pipelines detected.
          </Box>
        </Show>
        <Show when={pipelines.error != null}>
          <Box bg="$neutral7" borderRadius="$md" p="$3" fontStyle="italic" fontSize="$xs">
            <Box mt="$2">
              <Box fontSize="$xs">Error fetching pipeline information:</Box>
              <Box fontSize="$xs">by {pipelines.error}</Box>
            </Box>
          </Box>
        </Show>
      </Skeleton>
    </Box>
  )
}
