import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Controller, type SubmitHandler, useForm } from "react-hook-form";
import { type Review, saveReview, fetchMovie } from "@/services/api/movies";
import { Button } from "./ui/button";
import {
    Field,
    FieldError,
    FieldLabel,
    FieldSeparator,
    FieldSet,
} from "./ui/field";
import { Input } from "./ui/input";
import { Slider } from "./ui/slider";
import { Textarea } from "./ui/textarea";

export default function ReviewForm({
    review,
    movieId,
    onSuccess: userSucess,
    onCancel,
}: {
    onSuccess: () => void;
    onCancel?: () => void;
    review?: Review;
    movieId: number;
}) {
    const queryClient = useQueryClient();

    const { control, handleSubmit, reset } = useForm<Review>({
        mode: "onChange",
        defaultValues: {
            title: review?.title ?? "",
            rating: review?.rating ?? 5,
            review: review?.review ?? "",
        },
    });

    const resetForm = () => {
        reset({
            title: review?.title ?? "",
            rating: review?.rating ?? 5,
            review: review?.review ?? "",
        });
    };

    const mutation = useMutation({
        mutationFn: (data: Review) =>
            saveReview(movieId, review?.id ?? 0, data),
        onSuccess: () => {
            reset();
            queryClient.invalidateQueries({
                queryKey: ["movies", movieId],
            });

            queryClient.invalidateQueries({
                queryKey: ["movies", movieId, "reviews", "me"],
            });
            userSucess();
        },
    });

    const onSubmit: SubmitHandler<Review> = (data) => {
        mutation.mutate(data);
    };

    const MAX_TITLE = 100;
    const MAX_REVIEW = 1000;
    return (
        <form
            onSubmit={handleSubmit(onSubmit)}
            className="mt-6 flex flex-col gap-3"
        >
            <FieldSet>
                <Controller
                    name="title"
                    control={control}
                    rules={{
                        required: "Title required",
                        maxLength: {
                            value: MAX_TITLE,
                            message: `Max ${MAX_TITLE} chars`,
                        },
                    }}
                    render={({ field, fieldState }) => (
                        <Field data-invalid={fieldState.invalid}>
                            <div className="flex space-x-4">
                                <FieldLabel> Title </FieldLabel>
                                <Input
                                    {...field}
                                    id={field.name}
                                    aria-invalid={fieldState.invalid}
                                    placeholder="Amazing..."
                                    autoComplete="off"
                                />
                            </div>
                            {fieldState.invalid && (
                                <FieldError errors={[fieldState.error]} />
                            )}
                        </Field>
                    )}
                />
                <FieldSeparator />
                <Controller
                    name="rating"
                    control={control}
                    render={({ field, fieldState }) => (
                        <Field data-invalid={fieldState.invalid}>
                            <div className="flex space-x-4">
                                <FieldLabel>Rating</FieldLabel>
                                <Slider
                                    value={[field.value]}
                                    onValueChange={(val) =>
                                        field.onChange(val[0])
                                    }
                                    step={1}
                                    max={10}
                                    aria-label="Price Range"
                                />
                                <span>{field.value}</span>
                            </div>
                            {fieldState.invalid && (
                                <FieldError errors={[fieldState.error]} />
                            )}
                        </Field>
                    )}
                ></Controller>
                <FieldSeparator />
                <Controller
                    name="review"
                    control={control}
                    rules={{
                        required: "Review required",
                        maxLength: {
                            value: MAX_REVIEW,
                            message: `Max ${MAX_REVIEW} chars`,
                        },
                    }}
                    render={({ field, fieldState }) => (
                        <Field data-invalid={fieldState.invalid}>
                            <FieldLabel>Review</FieldLabel>
                            <Textarea
                                {...field}
                                placeholder="I liked this movie..."
                                aria-invalid={fieldState.invalid}
                                autoComplete="off"
                                ref={(el) => {
                                    field.ref(el);
                                    if (el)
                                        el.style.height = `${el.scrollHeight}px`;
                                }}
                                onInput={(e) => {
                                    e.currentTarget.style.height = "auto";
                                    e.currentTarget.style.height = `${e.currentTarget.scrollHeight}px`;
                                    field.onChange(e);
                                }}
                                className="max-h-96 min-h-56 overflow-hidden"
                            />
                            {fieldState.invalid && (
                                <FieldError errors={[fieldState.error]} />
                            )}
                        </Field>
                    )}
                />
                <Field className="justify-end" orientation="horizontal">
                    {mutation.isError && (
                        <p className="mr-auto text-red-500 text-sm">
                            {mutation.error?.message || "Something went wrong"}
                        </p>
                    )}
                    <Button type="submit" disabled={mutation.isPending}>
                        {mutation.isPending ? "Submitting..." : "Submit Review"}
                    </Button>

                    <Button type="button" variant="outline" onClick={resetForm}>
                        Reset
                    </Button>
                    {onCancel && (
                        <Button
                            type="button"
                            variant="outline"
                            onClick={onCancel}
                        >
                            Cancel
                        </Button>
                    )}
                </Field>
            </FieldSet>
        </form>
    );
}
