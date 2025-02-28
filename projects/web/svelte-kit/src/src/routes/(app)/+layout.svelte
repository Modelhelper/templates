<script lang="ts">
	import Header from "$lib/shared/Header.svelte";
  import { setDefaultGaugeState } from "$lib/state/metric-options.svelte";
  import { quintInOut, cubicIn, cubicOut } from "svelte/easing";
  import { fade, slide, fly } from "svelte/transition";
  export let data: { pathname: string, user: any };


  $: ({ pathname } = data);

  const duration = 300;
	const delay = duration + 100;
	const y = 10;

	const transitionIn = { easing: cubicOut, y, duration, delay };
	const transitionOut = { easing: cubicIn, y: -y, duration };

  setDefaultGaugeState();
</script>

<Header user={data.user} isAdmin={data.user && data.user.email.toString().includes('eitvet')}/>

{#key pathname}
  <main class="w-full" in:fly={transitionIn} out:fly={transitionOut} >
    <slot />
  </main>
{/key}

