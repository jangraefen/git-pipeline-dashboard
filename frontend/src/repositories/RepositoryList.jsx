import { createReaction, createResource, createSignal, For, onCleanup, Show } from "solid-js"
import { Box, Skeleton } from "@hope-ui/solid"
import RepositoryGroup from "./RepositoryGroup"

const repositoryFetcher = async () => (await fetch("/repositories")).json()

export default () => {
  const [repositoryGroups, { refetch: refetchRepositoryGroups }] = createResource(repositoryFetcher)

  const [initiallyLoaded, setInitiallyLoaded] = createSignal(false)
  const loadTracker = createReaction(() => setInitiallyLoaded(true))
  loadTracker(repositoryGroups)

  const interval = setInterval(refetchRepositoryGroups, 5 * 60 * 1000)
  onCleanup(() => clearInterval(interval))

  return (
    <Box w="100%">
      <Skeleton loaded={initiallyLoaded}>
        <Show when={repositoryGroups.error == null && repositoryGroups() != null && !repositoryGroups().status}>
          <For each={repositoryGroups()}>{(group) => <RepositoryGroup title={group.title} repositories={group.repositories} />}</For>
        </Show>
        <Show when={repositoryGroups.error == null && repositoryGroups() != null && !!repositoryGroups().status}>
          <Box bg="$neutral7" borderRadius="$md" p="$3" fontStyle="italic" fontSize="$xs">
            <Box mt="$2">
              <Box fontSize="$xs">Error while fetching repositories:</Box>
              <Box fontSize="$xs">
                {repositoryGroups().status} - {repositoryGroups().title}
              </Box>
              <Box fontSize="$xs">{repositoryGroups().detail}</Box>
            </Box>
          </Box>
        </Show>
        <Show when={repositoryGroups.error != null}>
          <Box bg="$neutral7" borderRadius="$md" p="$3" fontStyle="italic" fontSize="$xs">
            <Box mt="$2">
              <Box fontSize="$xs">Error fetching pipeline information:</Box>
              <Box fontSize="$xs">by {repositoryGroups.error}</Box>
            </Box>
          </Box>
        </Show>
      </Skeleton>
    </Box>
  )
}
