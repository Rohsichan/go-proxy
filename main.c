#if 0
#include <stdio.h>
#include <string.h>
#include <pcre.h>

int main(int argc, char **argv)
{
	pcre *re;
	int error_offset;
	const char *error;
	int ret;
	int len;
	int ovector_qp[3] = {0, };
	int i;

	char buf[1000];

	memset(buf, 0, sizeof(buf));

	char *pattern = "{{(?:\\n|\\r\\n?|.)?[^\\w\\d\\s-,:_]+(?:\\n|\\r\\n?|.)?}}";

	printf("==## pattern [%s]\n", pattern);

	re = pcre_compile(pattern, 1, &error, &error_offset, NULL);

	if (NULL == re) {
			printf("re is NULL\n");
			printf("==## erorr [%s], error_offset [%d]\n", error, error_offset);
		return -1;
	}


	u_int32_t matchlimit;
	ret = pcre_fullinfo(re, NULL, PCRE_INFO_MATCHLIMIT, &matchlimit);
#if 0
	if (ret != 0) {
			printf("pcre matchlimit error %d\n", ret);
			return -1;
	}

	printf("matchlimit %d\n", matchlimit);
return 0;
#endif

	FILE *ptr_file;

	ptr_file = fopen("./pattern_test", "r");
	if (NULL == ptr_file) {
			printf("ptr_file is null");
			return -1;
	}

	len = fread(buf, 1, 1000, ptr_file);
	fclose(ptr_file);


	printf("\n[");
	for (i = 0; i < len; i++) {
		printf("%c", buf[i]);
	}
	printf("]\n");

	printf("\n len %d\n", len);

	ret = pcre_exec(re, NULL, buf, len, 0, 0, ovector_qp, 3);

	printf("ret   [%d]\n", ret);


	return 0;
}
#else
#include <stdio.h>
#include <string.h>
#include <pcre.h>

int main() {
    const char* pattern = "^[a-zA-Z]+$";
    const char* input = "";
	int ovector[3] = {0, };
	int rc;

    pcre* re = pcre_compile(pattern, 0, NULL, NULL, NULL);
    rc = pcre_exec(re, NULL, input, 0, 0, ovector, 3);
    if (rc == -1) {
        printf("Matching error: %d\n", rc);
    } else if (rc == -8) {
        printf("Empty string not allowed.\n");
    } else {
        printf("Match found: %d\n", rc);
    }
    pcre_free(re);
    return 0;
}
#endif
