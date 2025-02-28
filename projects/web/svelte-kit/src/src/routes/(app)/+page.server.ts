import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ depends, locals: { supabase } }) => {
	depends('supabase:db:task');
	const { data: task } = await supabase.from('task')
        .select('id,title,task_status(title),task_status_id,task_type_id,task_type(title)')
		// .where('task_status_id', 'eq', 'new')
		.neq('task_status_id', 'done')
        .order('task_status_id', 'id', { ascending: false });

	// depends('supabase:db:preset');
	// const {preset} = await supabase.from('preset').select('id, title, body')
	
	// return { tasks: task ?? [], presets: preset ?? [] };
	return { tasks: task ?? [] };
};