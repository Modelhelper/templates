<script lang="ts">
	import pages from '$lib/stores/pages';
	import type { PageData } from './$types';
	// import type { PageData } from './$types';

	export let data: PageData;
</script>

<div class="p-4 px-16 flex gap-4 wrap justify-between">
	<div class="flex flex-col">
		<h1 class="text-4xl mb-2">
			Welcome to this Configurator | Builder thingy <small class="text-sm">Alpha</small>
		</h1>

		{#if !data.user}
			<p class="text-lg">You are not logged in</p>
			<div>
				If you have an account, please go to the <span class="text-blue-400"
					><a href="/login">login page</a></span
				>
				to log in or create a user

				<div
					class="mt-4 p-4 mb-4 text-sm text-red-800 rounded-lg bg-red-50 dark:bg-gray-800 dark:text-red-400"
					role="alert"
				>
					<p class="mt-2">
						Access to these tools are by <span class="font-bold">invitation</span> only, please contact
						me if you want to use this service (you probably know who I am...)
					</p>
				</div>
			</div>
		{:else}
			<p class="text-lg">Hello, {data.user.email}</p>

			<p>
				This project is built with Svelte 5 (in beta) and SvelteKit (v2). For design I have used
				TailwindCSS 4 (in alpha|beta)
			</p>

			<h2 class="mt-8"><b>A few things</b></h2>
			<ul class="list-disc pl-6">
				<li>You cannot save a preset</li>
				<li>the state is not shared between windows/tabs</li>
				<li>
					there are a few glitches with some of the settings (ex: Active Segment.FillColor =
					Gradient)
				</li>
				<li>You cannot add or remove sensors/riders</li>
				<li>It's not possible to randomize sensor metrics</li>
			</ul>

			<h2 class="my-4"><b>Pages</b></h2>
			<div class="flex flex-row flex-wrap gap-4">
				{#each $pages as page}
					<a href={page.path} class="bg-blue-400 text-white p-2 rounded-lg">
						<section class="w-80 h-48">
							<h3 class="text-2xl font-bold">{page.name}</h3>
							<p>{page.description}</p>
						</section>
					</a>
				{/each}
			</div>

			<h2 class="mt-8"><b>About the sensor data:</b></h2>
			<p>
				Initially there are 8 sensors. They are set up so that they represent a unique zone (in
				order). The last one is way above the max.
			</p>
			<p>
				The "maxed-out" sensor is there so you can play with the 'At Threshold' event. E.g what will
				happen with the metric when a sensor goes above the threshold (160% (W/ftp))
			</p>
		{/if}
	</div>

	{#if data.user}
		<div class="min-w-96 flex flex-col gap-4">
			<section>
				<h2 class="my-3 text-2xl font-bold">Some of the stuff to do...</h2>
				<ul class="list-disc pl-6">
					{#each data.tasks as task}
						<li class="flex justify-between my-2">
							<span>
								<span class:epic-badge={task.task_type_id==='epic'} class:task-badge={task.task_type_id==='task'} >
									{task.task_type.title} 
								</span>
								&nbsp;
								<!-- class:decoration-pink-500={task.task_status_id === 'done'} -->
								<span class:line-through={task.task_status_id === 'done'} >{task.title}</span>
							 </span>
							<span class:new-badge={task.task_status_id==='new'} class:done-badge={task.task_status_id==='done'} class:inprogress-badge={task.task_status_id==='inprogress'}>{task.task_status.title}</span>
						</li>
					{/each}
					
				</ul>
			</section>
		</div>
	{/if}
</div>
		<!-- <pre>{JSON.stringify(data.presets, null, 2)}</pre> -->
