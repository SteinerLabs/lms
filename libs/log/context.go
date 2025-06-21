package log

import "context"

// ContextWithValue adds a key-value pair to the context's logWithFields map.
// If the map is nil, it creates a new map and assigns it to the context.
// If the map already exists, it performs a shallow copy of the existing map,
// adds the key-value pair to the copy, and assigns it to the context.
// Lastly, it returns the updated context.
// This function is used to attach additional log fields to a specific context.
func ContextWithValue(ctx context.Context, key string, value any) context.Context {
	ctxValue := ctx.Value(logWithFields)
	if ctxValue == nil {
		return context.WithValue(ctx, logWithFields, map[string]any{key: value})
	}

	fields := shallowCopy(ctxValue.(map[string]interface{}))
	fields[key] = value

	return context.WithValue(ctx, logWithFields, fields)
}

// ContextWithValues creates a new context with additional key-value pairs specified in the fields parameter.
// If the supplied context already contains a value for the key "logWithFields", then the function merges the existing fields
// with the new fields using the "merge" function. Otherwise, a new logWithFields key-value pair is added to the context
// using the shallowCopy of the fields.
func ContextWithValues(ctx context.Context, fields map[string]any) context.Context {
	ctxValue := ctx.Value(logWithFields)

	if ctxValue == nil {
		return context.WithValue(ctx, logWithFields, shallowCopy(fields))
	}

	return context.WithValue(ctx, logWithFields, merge(ctxValue.(map[string]interface{}), fields))
}
