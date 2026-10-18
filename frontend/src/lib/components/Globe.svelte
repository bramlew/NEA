<script lang="ts">

    import { onDestroy, onMount } from "svelte"
    import Globe, { type GlobeInstance } from "globe.gl"
    import { decode } from "@googlemaps/polyline-codec"
    import type { Response } from "#lib/types.ts";

    let { routeData } : { routeData: Response } = $props()

    let globeContainer: HTMLElement
    let globe: GlobeInstance
    let roadRoutes = $derived(routeData.details.route.filter(leg => leg.isRoad))
    let roadRoutesPaths = $derived(roadRoutes.map(leg => decode(leg.polyline)))
    let airRoutes = $derived(routeData.details.route.filter(leg => !leg.isRoad))
    let airRoutesArcs = $derived(airRoutes.map(leg => { return {
        startLat: leg.origin.lat,
        startLng: leg.origin.lon,
        endLat: leg.dest.lat,
        endLng: leg.dest.lon
    }}))
    let width = $state(0)
    let height = $state(0)

    onMount(() => {
        globe = new Globe(globeContainer)
            .globeTileEngineUrl((x, y, l) => `https://tiles.versatiles.org/tiles/satellite/${l}/${x}/${y}`)
            .backgroundColor("#091224")
    })

    onDestroy(() => globe._destructor())

    $effect(() => {
        globe.width(width).height(height)
    })

    $effect(() => {
        globe.pathsData(roadRoutesPaths)
            .arcsData(airRoutesArcs)
            .arcAltitudeAutoScale(0.3)
    })

</script>

<div class="flex-2 w-1/2">
    <div bind:this={globeContainer} bind:clientWidth={width} bind:clientHeight={height} ></div>
</div>
