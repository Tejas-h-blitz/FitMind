import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { supabase } from './supabase';

export interface APIResponse<T> {
	success: boolean;
	data?: T;
	error?: string;
}

export interface UserProfile {
	id: string;
	full_name: string;
	streak_count: number;
	last_active: string;
	created_at: string;
}

export interface Document {
	id: string;
	user_id: string;
	name: string;
	size: number;
	storage_path: string;
	status: 'pending' | 'processing' | 'ready' | 'failed';
	created_at: string;
}

export interface SourceChunk {
	text: string;
	page_number: number;
}

export interface Message {
	id: string;
	doc_id: string;
	user_id: string;
	role: 'user' | 'assistant';
	content: string;
	sources: SourceChunk[];
	created_at: string;
}

export interface HealthMetric {
	id: string;
	user_id: string;
	bmi: number;
	height: number;
	weight: number;
	recorded_at: string;
}

export interface Goal {
	id: string;
	user_id: string;
	title: string;
	target_date: string;
	status: 'active' | 'completed';
	created_at: string;
}

async function getHeaders(isMultipart = false): Promise<HeadersInit> {
	const { data } = await supabase.auth.getSession();
	const token = data.session?.access_token || '';

	const headers: Record<string, string> = {
		'Authorization': `Bearer ${token}`
	};

	if (!isMultipart) {
		headers['Content-Type'] = 'application/json';
	}

	return headers;
}

async function request<T>(
	method: string,
	path: string,
	body?: any,
	isMultipart = false
): Promise<APIResponse<T>> {
	const url = `${PUBLIC_API_BASE_URL}${path}`;
	const headers = await getHeaders(isMultipart);

	const config: RequestInit = {
		method,
		headers
	};

	if (body !== undefined) {
		config.body = isMultipart ? body : JSON.stringify(body);
	}

	try {
		const res = await fetch(url, config);
		if (res.status === 401) {
			return { success: false, error: 'Session expired. Please log in again.' };
		}

		// SSE streams or other endpoints might not return JSON, but our standard REST handlers always do.
		const data = await res.json();
		return data as APIResponse<T>;
	} catch (err: any) {
		return { success: false, error: err.message || 'Network request failed' };
	}
}

export const api = {
	// Auth
	verifyAuth: () => request<UserProfile>('POST', '/api/auth/verify'),

	// Documents
	listDocuments: () => request<Document[]>('GET', '/api/documents'),
	
	uploadDocument: (file: File) => {
		const formData = new FormData();
		formData.append('file', file);
		return request<Document>('POST', '/api/documents/upload', formData, true);
	},
	
	deleteDocument: (id: string) => request<{ message: string }>('DELETE', `/api/documents/${id}`),
	
	getDocumentStatus: (id: string) => request<{ id: string; status: Document['status'] }>('GET', `/api/documents/${id}/status`),

	// Chat
	getChatHistory: (docId: string) => request<Message[]>('GET', `/api/chat/${docId}/history`),
	
	clearChatHistory: (docId: string) => request<{ message: string }>('DELETE', `/api/chat/${docId}/history`),

	// Health Tracking
	getHealthMetrics: () => request<HealthMetric[]>('GET', '/api/health/metrics'),
	
	createHealthMetric: (height: number, weight: number) => 
		request<HealthMetric>('POST', '/api/health/metrics', { height, weight }),
		
	getGoals: () => request<Goal[]>('GET', '/api/health/goals'),
	
	createGoal: (title: string, targetDate: string) => 
		request<Goal>('POST', '/api/health/goals', { title, target_date: targetDate }),
		
	updateGoalStatus: (id: string, status: 'active' | 'completed') => 
		request<{ id: string; status: string }>('PATCH', `/api/health/goals/${id}`, { status })
};
