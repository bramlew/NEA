<script lang="ts">
    import LocationInput from "#lib/components/LocationInput.svelte";
    import Globe from "#lib/components/Globe.svelte";
    import StatsSidebar from "#lib/components/StatsSidebar.svelte";
    import type { KeyedLeg, Response } from "#lib/types.ts"

    let routeData: Response = $state({ details: { route: [], totalWeight: 0 } })
    let legs: KeyedLeg[] = $derived(routeData.details.route.map((leg, i) => { return { id: i+1, details: leg} }))
    let totalWeight = $derived(routeData.details.totalWeight)

    const initRouteData = (data: Response) => {
        routeData = data
    }
</script>

<div class="flex">
    <LocationInput { initRouteData } />
    <Globe { legs } />
    <StatsSidebar { legs } { totalWeight } />
</div>
