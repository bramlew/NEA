<script lang="ts">
    import LocationInput from "#lib/components/LocationInput.svelte";
    import Globe from "#lib/components/Globe.svelte";
    import StatsSidebar from "#lib/components/StatsSidebar.svelte";
    import type { Response } from "#lib/types.ts"

    let routeData: Response = $state({ details: { route: [], totalWeight: 0 } })

    const initRouteData = (data: Response) => {
        routeData = data
    }
</script>

<div class="mx-auto flex">
    <LocationInput { initRouteData } />
    <Globe {routeData} />
    <StatsSidebar { routeData } />
</div>
