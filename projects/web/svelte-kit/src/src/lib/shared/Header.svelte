<script>
	let segmentId = $state(1);

	let { user, isAdmin = false } = $props();

	// let color = $derived(
	//     [segmentId],
	//     ([$segmentId]) => {
	//         return $segmentId === 1 ? 'red' : 'blue';
	//     }
	// );
</script>

<nav class="">
	<div class="left">
		<a href="/">
			<h1 class="text-2xl font-extrabold">Configurators</h1>
		</a>
	</div>

	<div class="right flex flex-row gap-4 font-bold">
		{#if user}
			<!-- <a href="/">Home</a> | -->
			<a href="/metric">Individual Metric</a>
			<span class="font-light text-neutral-400">|</span>
			<a href="/sensors">Sensors</a>
			<span class="font-light text-neutral-400">|</span>
			<a href="/states">State</a>

			{#if isAdmin}
				<span class="font-light text-neutral-400">|</span>
			    <a href="/plans">Plans</a>
				<span class="font-light text-neutral-400">|</span>
			    <a href="/hpe">WiP</a>
				<span class="font-light text-neutral-400">|</span>
			    <a href="/hpe">Admin</a>
			{/if}

			<span class="font-light text-neutral-400">|</span>
			<form action="/logout" method="POST">
				<button class="cursor-pointer" type="submit">Logout</button>
				<!-- <a type href="/logout">Logout</a>  -->
			</form>
		{:else}
			<a href="/login">Login or Sign up</a>
		{/if}
		
	</div>
</nav>
