<script lang="ts">
    import { Accordion } from "@skeletonlabs/skeleton-svelte";
    import type { KeyedLeg } from "#lib/types.ts"

    // Obtain the legs requested
    let { legs, totalWeight } : { legs: KeyedLeg[], totalWeight: number } = $props()
    
</script>

<div class="flex-1">
    <h3 class="h3">Route stats</h3>
    {#if legs.length > 0}
        <h6 class="h6">Total route length: {totalWeight}</h6>
        <Accordion collapsible>
            {#each legs as leg (leg.id)}
                <Accordion.Item value={leg.id.toString()}>
                    <h5 class="h5">
                        <Accordion.ItemTrigger>Leg {leg.id}</Accordion.ItemTrigger>
                    </h5>
                    <Accordion.ItemContent>
                        <p class="text-base">Origin coords: {leg.details.origin.lon}, {leg.details.origin.lat}</p>
                        <p class="text-base">Destination coords: {leg.details.dest.lon}, {leg.details.dest.lat}</p>
                        <p class="text-base">Distance: {leg.details.distance}</p>
                        <p class="text-base">Is a road: {leg.details.isRoad ? "Yes" : "No"}</p>
                    </Accordion.ItemContent>
                </Accordion.Item>
            {/each}
        </Accordion>
    {:else}
        <p class="text-base">Awaiting location request</p>
    {/if}
</div>
