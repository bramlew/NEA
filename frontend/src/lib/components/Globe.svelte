<script lang="ts">

    import { onDestroy, onMount } from "svelte"
    import Globe, { type GlobeInstance } from "globe.gl"
    import { decode } from "@googlemaps/polyline-codec"
    import type { Leg, Response } from "#lib/types.ts";

    let { routeData } : { routeData: Response } = $props()

    let globeContainer: HTMLElement
    let globe: GlobeInstance

    onMount(() => {
        let roadRoutes = $derived(routeData.details.route.filter(leg => leg.isRoad))
        let roadRoutesPaths = $derived(roadRoutes.map(leg => decode(leg.polyline)))
        let airRoutes = $derived(routeData.details.route.filter(leg => !leg.isRoad))

        const globe = new Globe(globeContainer)
            .globeTileEngineUrl((x, y, l) => `https://tiles.versatiles.org/tiles/satellite/${l}/${x}/${y}`)
            .pathsData(roadRoutesPaths)
    })

</script>

<div class="mx-auto flex-2">
    <div class="mx-auto" bind:this={globeContainer}></div>
</div>
