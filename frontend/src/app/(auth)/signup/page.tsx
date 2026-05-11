"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useSetAtom } from "jotai";
import { Box, Button, Heading, Input, Text, VStack } from "@chakra-ui/react";
import { useSignup } from "@/api/generated/ninjo";
import { userAtom } from "@/stores/auth";
import Link from "next/link";

export default function SignupPage() {
  const router = useRouter();
  const setUser = useSetAtom(userAtom);
  const signupMutation = useSignup();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (password.length < 8) {
      setError("パスワードは8文字以上で入力してください");
      return;
    }

    signupMutation.mutate(
      { data: { email, password, display_name: displayName } },
      {
        onSuccess: (user) => {
          setUser(user);
          router.push("/");
        },
        onError: (err) => {
          setError(err instanceof Error ? err.message : "登録に失敗しました");
        },
      }
    );
  };

  return (
    <Box as="form" onSubmit={handleSubmit} w="full" bg="white" p={8} borderRadius="lg" shadow="sm">
      <VStack gap={4}>
        <Heading size="lg">新規登録</Heading>

        {error && (
          <Text color="red.500" fontSize="sm">
            {error}
          </Text>
        )}

        <VStack gap={3} w="full">
          <Input
            type="text"
            placeholder="表示名"
            value={displayName}
            onChange={(e) => setDisplayName(e.target.value)}
            required
            maxLength={100}
          />
          <Input
            type="email"
            placeholder="メールアドレス"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <Input
            type="password"
            placeholder="パスワード（8文字以上）"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
          />
        </VStack>

        <Button
          type="submit"
          w="full"
          colorPalette="blue"
          loading={signupMutation.isPending}
        >
          登録する
        </Button>

        <Text fontSize="sm" color="gray.600">
          すでにアカウントをお持ちの方は{" "}
          <Link href="/login" style={{ color: "var(--chakra-colors-blue-500)" }}>
            ログイン
          </Link>
        </Text>
      </VStack>
    </Box>
  );
}
