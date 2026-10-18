<script lang="ts">
    import { Accordion } from "@skeletonlabs/skeleton-svelte";
    import type { Response, Leg, KeyedLeg } from "#lib/types.ts"

    // Obtain the legs requested
    let { routeData } : { routeData: Response } = $props()
    let legs: KeyedLeg[] = $derived(routeData.details.route.map((leg: Leg, i: number) => { return { id: i+1, leg: leg} }))
</script>

<div class="flex-1">
    <h3 class="h3">Some stats</h3>
    <Accordion collapsible>
        {#each legs as leg (leg.id)}
            <Accordion.Item value={leg.id.toString()}>
                <h5 class="h5">
                    <Accordion.ItemTrigger>Leg {leg.id}</Accordion.ItemTrigger>
                </h5>
                <Accordion.ItemContent>
                    <p class="text-base">Origin coords: {leg.leg.origin.lon}, {leg.leg.origin.lat}</p>
                    <p class="text-base">Destination coords: {leg.leg.dest.lon}, {leg.leg.dest.lat}</p>
                    <p class="text-base">Distance: {leg.leg.distance}</p>
                    <p class="text-base">Is a road: {leg.leg.isRoad ? "Yes" : "No"}</p>
                </Accordion.ItemContent>
            </Accordion.Item>
        {/each}
    </Accordion>
</div>
